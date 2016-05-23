
#include <TinyGPS++.h>

#include <RH_RF95.h>

/* for feather32u4 */
#define RFM95_CS 8
#define RFM95_RST 4
#define RFM95_INT 7

#define RF95_FREQ 915.0

RH_RF95 rf95(RFM95_CS, RFM95_INT);

TinyGPSPlus gps;

#define LED 13

void setup() 
{
  Serial.begin(9600);
  Serial1.begin(9600);

  digitalWrite(4, HIGH);
  pinMode(LED, OUTPUT);
  digitalWrite(LED, HIGH);

  if (!rf95.init())
    error(1);

  if (!rf95.setFrequency(RF95_FREQ))
    error(2);

  rf95.setModemConfig(RH_RF95::Bw31_25Cr48Sf512);

  rf95.setTxPower(20);
}

struct RecBuffer {
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
} recBuffer;

uint8_t buf[RH_RF95_MAX_MESSAGE_LEN];
uint8_t len = sizeof(buf);

void loop() {
  //feed the ever-hungry gps
  while (Serial1.available() > 0)
    gps.encode(Serial1.read());

  if (rf95.waitAvailableTimeout(800)) {
    // Should be a reply message for us now   
    if (rf95.recv(buf, &len)) {

      memcpy(&recBuffer, buf, sizeof(recBuffer));

      // write the sentence to serial
      Serial.print(rf95.lastRssi(), DEC);
      Serial.print(",");
      Serial.print(recBuffer.count, DEC);
      Serial.print(",");
      Serial.print(recBuffer.raw_lat_deg + recBuffer.raw_lat_bil / 1000000000.0, 6);
      Serial.print(",-"); // xxx: bad - hard coded negative here
      Serial.print(recBuffer.raw_lng_deg + recBuffer.raw_lng_bil / 1000000000.0, 6);
      Serial.print(",");
      Serial.print(recBuffer.alt / 100.0);
      Serial.print(",");
      Serial.print(recBuffer.speed / 100.0);
      Serial.print(",");
      Serial.print(recBuffer.tmpin / 100.0);
      Serial.print(",");
      Serial.print(recBuffer.tmpout / 100.0);
      Serial.print(",");
      Serial.print(recBuffer.press / 100.0);
      Serial.print(",");
      Serial.print(recBuffer.voltage / 10.0);
      Serial.print(",");
      Serial.print(gps.location.lat(), 6);
      Serial.print(",");
      Serial.println(gps.location.lng(), 6);

      blink_led();
    }
  }
}

void blink_led() {
  digitalWrite(LED, HIGH);
  delay(100);
  digitalWrite(LED, LOW);
}

// blink out an error code
void error(uint8_t errno) {
  while(1) {
    uint8_t i;
    for (i=0; i<errno; i++) {
      digitalWrite(LED, HIGH);
      delay(100);
      digitalWrite(LED, LOW);
      delay(300);
    }
    delay(1000);
  }
}

