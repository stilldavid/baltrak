var map = L.map('map').setView([40.035, -105.269], 10);

L.tileLayer('http://localhost:3000/tiles/{z}/{x}/{y}.png', {
  attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="http://mapbox.com">Mapbox</a>',
}).addTo(map);


var latlngs = [];
var polyline = L.polyline(latlngs, {color: 'red'}).addTo(map);


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
}, 1000);

// set up a web socket
var serversocket = new WebSocket("ws://localhost:3000/ws");

// Write message on receive
serversocket.onmessage = function(e) {
  var obj = jQuery.parseJSON(e.data);
  polyline.addLatLng([obj.lat, obj.lng]);

  $('#rssi').text(obj.rssi);
  $('#count').text(obj.count);
  $('#lat').text(obj.lat);
  $('#lng').text(obj.lng);
  $('#alt').text(obj.alt);
  $('#spd').text(obj.spd);
  $('#tmpint').text(obj.tmpint);
  $('#tmpext').text(obj.tmpext);
  $('#volts').text(obj.volts);

  // update the timer
  lastUpdate = moment();
};

