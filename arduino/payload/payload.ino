
#include <TinyGPS++.h>
#include <SPI.h>
#include <RH_RF95.h>

// BME280
#include <Wire.h>
#include <SPI.h>
#include <Adafruit_BME280.h>
#include <Adafruit_MCP9808.h>

// writing to flash SD card
#include <SD.h>

#define SEALEVELPRESSURE_HPA (1013.25)

// to measure battery voltage
#define VBATPIN A7

// Set the pins used
#define cardSelect 4
File logfile;
 
/* for feather32u4 */
#define RFM95_CS 6
#define RFM95_RST 9
#define RFM95_INT 5

// Change to 434.0 or other frequency, must match RX's freq!
#define RF95_FREQ 915.0
 
// Singleton instance of the radio driver
RH_RF95 rf95(RFM95_CS, RFM95_INT);
 
// Blinky on receipt
#define LED 13

TinyGPSPlus gps;

// internal BME280 temp/pressure/humidity
Adafruit_BME280 bme; // I2C

// external MCP9808 temperature sensor
Adafruit_MCP9808 tempsensor = Adafruit_MCP9808();

void setup() {
  Serial.begin(9600);
  Serial1.begin(9600);

  pinMode(LED, OUTPUT);
  pinMode(8, OUTPUT);
  digitalWrite(8, HIGH);

  if (!rf95.init())
    error(1);

  if (!rf95.setFrequency(RF95_FREQ))
    error(1);

  rf95.setModemConfig(RH_RF95::Bw31_25Cr48Sf512);

  rf95.setTxPower(20);

  // bme280 - internal temp sensor
  if (! bme.begin()) {
    error(2);
  }

  // mcp9808 - external temp
  if (! tempsensor.begin()) {
    error(3);
  }

  // see if the card is present and can be initialized:
  if (! SD.begin(cardSelect)) {
    error(4);
  }

  char filename[15];
  strcpy(filename, "BAL00.TXT");
  for (uint8_t i = 0; i < 100; i++) {
    filename[3] = '0' + i/10;
    filename[4] = '0' + i%10;
    // create if does not exist, do not open existing, write, sync after write
    if (! SD.exists(filename)) {
      break;
    }
  }

  logfile = SD.open(filename, FILE_WRITE);
  if(! logfile ) {
    error(5);
  }

  // get some data in gps buffer
  delay(1000);
}

// data to send - this _must_ match up with the struct on the recieve side
struct SendBuffer {
  uint32_t count;
  uint32_t raw_lat_deg;
  uint32_t raw_lat_bil;
  uint32_t raw_lng_deg;
  uint32_t raw_lng_bil;
  int32_t alt;
  int32_t speed;
  int32_t tmpin;
  int32_t tmpout;
  int32_t press;
  int32_t voltage;
} sendBuffer;

int next = 0;

// iterator
uint16_t counter = 0;

void loop() {

  // feed the ever-hungry gps
  while (Serial1.available() > 0)
    gps.encode(Serial1.read());

  sendBuffer.count = counter;

  // pull in temp data from i2c sensors
  tempsensor.shutdown_wake(0); // wake up external tmp
  sendBuffer.tmpout = (int32_t) (tempsensor.readTempC() * 100);
  //tempsensor.shutdown_wake(1); // sleep external tmp

  // measure the battery
  float measuredvbat = analogRead(VBATPIN);
  measuredvbat *= 2;    // we divided by 2, so multiply back
  measuredvbat *= 3.3;  // Multiply by 3.3V, our reference voltage
  measuredvbat /= 1024; // convert to voltage
  sendBuffer.voltage = (int32_t) (measuredvbat * 10);

  // pull data from gps into struct to send
  sendBuffer.raw_lat_deg = gps.location.rawLat().deg;
  sendBuffer.raw_lat_bil = gps.location.rawLat().billionths;
  sendBuffer.raw_lng_deg = gps.location.rawLng().deg;
  sendBuffer.raw_lng_bil = gps.location.rawLng().billionths;
  sendBuffer.alt = gps.altitude.value();
  sendBuffer.speed = gps.speed.value();

  sendBuffer.tmpin = (int32_t) (bme.readTemperature() * 100);
  sendBuffer.press = (int32_t) (bme.readPressure() * 100);

  // create an unsigned int array to hold the transmission message
  // and copy byte for byte into the array
  uint8_t b[sizeof(sendBuffer)];
  memcpy(b, &sendBuffer, sizeof(sendBuffer));

  // send the data, wait for the Tx to go
  rf95.send(b, sizeof(b));
  rf95.waitPacketSent();

  digitalWrite(8, HIGH);
  logfile.print(counter);
  logfile.print(",");
  logfile.print(gps.date.value());
  logfile.print(",");
  logfile.print(gps.time.value());
  logfile.print(",");
  logfile.print(sendBuffer.raw_lat_deg);
  logfile.print(",");
  logfile.print(sendBuffer.raw_lng_deg);
  logfile.print(",");
  logfile.print(sendBuffer.alt);
  logfile.print(",");
  logfile.print(sendBuffer.speed);
  logfile.print(",");
  logfile.print(sendBuffer.tmpin);
  logfile.print(",");
  logfile.print(sendBuffer.tmpout);
  logfile.print(",");
  logfile.print(sendBuffer.press);
  logfile.print(",");
  logfile.println(sendBuffer.voltage);

  digitalWrite(8, LOW);

  // flashythingie!
  blinkLed();

  counter++;

  // arbitrary amonut of delay
  smartDelay(850);
}

// This custom version of delay() ensures that the gps object
// is being "fed".
static void smartDelay(unsigned long ms)
{
  unsigned long start = millis();
  do 
  {
    while (Serial1.available())
      gps.encode(Serial1.read());
  } while (millis() - start < ms);
}

// happy little blinkies
void blinkLed() {
  digitalWrite(LED, HIGH);
  delay(100);
  digitalWrite(LED, LOW);
}

// blink out an error code... forever
void error(uint8_t errno) {
  while(1) {
    uint8_t i;
    for (i=0; i<errno; i++) {
      digitalWrite(13, HIGH);
      delay(100);
      digitalWrite(13, LOW);
      delay(300);
    }
    delay(1000);
  }
}

