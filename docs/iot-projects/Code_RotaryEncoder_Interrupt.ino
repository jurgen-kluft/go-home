/*
 * This ESP32 code is created by esp32io.com
 *
 * This ESP32 code is released in the public domain
 *
 * For more detail (instruction and wiring diagram), visit https://esp32io.com/tutorials/esp32-rotary-encoder
 */

#include <ezButton.h>  // the library to use for SW pin

#define CLK_PIN 25  // ESP32 pin GPIO25 connected to the rotary encoder's CLK pin
#define DT_PIN  26  // ESP32 pin GPIO26 connected to the rotary encoder's DT pin
#define SW_PIN  27  // ESP32 pin GPIO27 connected to the rotary encoder's SW pin

#define DIRECTION_CW  0  // clockwise direction
#define DIRECTION_CCW 1  // counter-clockwise direction

volatile int           counter   = 0;
volatile int           direction = DIRECTION_CW;
volatile unsigned long last_time;  // for debouncing
int                    prev_counter;

ezButton button(SW_PIN);  // create ezButton object that attach to pin 7;

void IRAM_ATTR ISR_encoder()
{
    if ((millis() - last_time) < 50)  // debounce time is 50ms
        return;

    if (digitalRead(DT_PIN) == HIGH)
    {
        // the encoder is rotating in counter-clockwise direction => decrease the counter
        counter--;
        direction = DIRECTION_CCW;
    }
    else
    {
        // the encoder is rotating in clockwise direction => increase the counter
        counter++;
        direction = DIRECTION_CW;
    }

    last_time = millis();
}

void setup()
{
    Serial.begin(9600);

    // configure encoder pins as inputs
    pinMode(CLK_PIN, INPUT);
    pinMode(DT_PIN, INPUT);
    button.setDebounceTime(50);  // set debounce time to 50 milliseconds

    // use interrupt for CLK pin is enough
    // call ISR_encoder() when CLK pin changes from LOW to HIGH
    attachInterrupt(digitalPinToInterrupt(CLK_PIN), ISR_encoder, RISING);
}

void loop()
{
    button.loop();  // MUST call the loop() function first

    if (prev_counter != counter)
    {
        Serial.print("Rotary Encoder:: direction: ");
        if (direction == DIRECTION_CW)
            Serial.print("CLOCKWISE");
        else
            Serial.print("ANTICLOCKWISE");

        Serial.print(" - count: ");
        Serial.println(counter);

        prev_counter = counter;
    }

    if (button.isPressed())
    {
        Serial.println("The button is pressed");
    }

    // TO DO: your other work here
}
