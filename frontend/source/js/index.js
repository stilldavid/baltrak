L.Icon.Default.imagePath = '/assets/img/';

var balloonIcon = L.icon({
    iconUrl: '/assets/img/balloon-red.png',
    iconSize:     [46, 84],
    iconAnchor:   [23, 84],
});

var chaseIcon = L.icon({
    iconUrl: '/assets/img/car-blue.png',
    iconSize:     [55, 25],
    iconAnchor:   [26, 25],
});

var map = L.map('map').setView([40.035, -105.269], 10);

L.tileLayer('http://localhost:3000/tiles/{z}/{x}/{y}.png', {
  attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="http://mapbox.com">Mapbox</a>',
}).addTo(map);

var latlngs = [];
var polyline = L.polyline(latlngs, {color: 'red'}).addTo(map);

var balloon = L.marker([50.5, 30.5], {icon: balloonIcon}).addTo(map);
var chase = L.marker([50.5, 30.5], {icon: chaseIcon}).addTo(map);

// zoom the map to the polyline
//map.fitBounds(polyline.getBounds());

// leaflet.coordinates plugin...
L.control.coordinates({
  position: "bottomright", //optional default "bootomright"
  decimals: 5, //optional default 4
  decimalSeperator: ".", //optional default "."
  labelTemplateLat: "Latitude: {y}", //optional default "Lat: {y}"
  labelTemplateLng: "Longitude: {x}", //optional default "Lng: {x}"
  enableUserInput: true, //optional default true
  useDMS: false, //optional default false
  useLatLngOrder: true, //ordering of labels, default false-> lng-lat
  markerType: L.marker, //optional default L.marker
  markerProps: {} //optional default {}
}).addTo(map);

// timer thing
var lastUpdate = moment();
setInterval(function() {
  var ago = moment().diff(lastUpdate, 'seconds');
  $('#lastupdate').text(ago);
}, 100);

// set up a web socket
var serversocket = new WebSocket("ws://localhost:3000/ws");

// Write message on receive
serversocket.onmessage = function(e) {
  var obj = jQuery.parseJSON(e.data);

  polyline.addLatLng([obj.lat, obj.lng]);
  balloon.setLatLng([obj.lat, obj.lng]);
  chase.setLatLng([obj.chase_lat, obj.chase_lng]);

  $('#rssi').text(obj.rssi);
  $('#count').text(obj.count);
  $('#lat').text(obj.lat);
  $('#lng').text(obj.lng);
  $('#alt').text(obj.alt);
  $('#spd').text(obj.spd);
  $('#spdmph').text((obj.spd * 2.236).toFixed(1));
  $('#tmpint').text(obj.tmpint);
  $('#tmpext').text(obj.tmpext);
  $('#press').text(obj.press);
  $('#volts').text(obj.volts);

  // update the timer
  lastUpdate = moment();
};

// pull in history
$.getJSON('/history.json', function(data) {
  $.each(data.history, function( key, val ) {
    if(0 !=val.lag && 0 != val.lng) {
      polyline.addLatLng([val.lat, val.lng]);
    }
  });
});
