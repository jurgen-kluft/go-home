#include <ESP8266WiFi.h>
#include <WiFiUdp.h>
#include "EEPROM.h"

const char*  ssid          = "OBNOSIS8";
const char*  password      = "abcdefghijkl8";
unsigned int udpLocalPort  = 31337;
unsigned int udpRemotePort = 31370;

const int reedSwitch = 13;
const int powerOff   = 16;  // set to low to turn off LDO

struct wifi_cache_t
{
    uint32_t crc;
    uint32_t ip_address;
    uint32_t ip_gateway;
    uint32_t ip_mask;
    uint32_t ip_dns1;
    uint32_t ip_dns2;
    uint16_t channel;
    uint8_t  bssid[6];
    void     reset();
    uint32_t calc_crc() const;
    bool     load();
    void     save();
    void     connect(const char* ssid, const char* auth);
};

void eeprom_save(byte const* data, int32_t size)
{
    EEPROM.begin(512);
    for (int32_t i = 0; i < size; i++)
        EEPROM.write(i, data[i]);
    EEPROM.commit();
}

void eeprom_load(uint8_t* data, int32_t size)
{
    EEPROM.begin(512);
    for (int32_t i = 0; i < size; i++)
        data[i] = EEPROM.read(i);
}

void wifi_cache_t::reset()
{
    crc        = 0;
    ip_address = 0;
    ip_gateway = 0;
    ip_mask    = 0;
    ip_dns1    = 0;
    ip_dns2    = 0;
    channel    = 0;
    // skip bssid
}

// Note: this is not a real CRC computation
uint32_t wifi_cache_t::calc_crc() const 
{
    uint16_t const* data = (uint16_t const*)this;
    uint16_t const* end  = (uint16_t const*)(this + 1);
    uint32_t c = 0;
    uint32_t m = 1;
    while (data < end)
    {
        c = c + (*data++ * (2*m + 1));
        m++;
    }
    return c;
}

bool wifi_cache_t::load()
{
    //eeprom_load((uint8_t*)this, sizeof(wifi_cache_t));
    //uint32_t loaded_crc = crc;
    //crc                 = 0;
    //uint32_t verify_crc = calc_crc();
    //if (loaded_crc != verify_crc)
    //{
    //     reset();
    //     return false;
    // }
    crc = 0;
    ip_address = uint32_t(236) << 24 | uint32_t(8) << 16   | uint32_t(168) << 8 | uint32_t(192);
    ip_gateway = uint32_t(1) << 24   | uint32_t(8) << 16   | uint32_t(168) << 8 | uint32_t(192);
    ip_mask    = uint32_t(0) << 24   | uint32_t(255) << 16 | uint32_t(255) << 8 | uint32_t(255);
    ip_dns1    = uint32_t(1) << 24   | uint32_t(8) << 16   | uint32_t(168) << 8 | uint32_t(192);
    ip_dns2    = 0;
    channel    = 6;
    bssid[0]   = 98;
    bssid[1]   = 55;
    bssid[2]   = 240;
    bssid[3]   = 134;
    bssid[4]   = 57;
    bssid[5]   = 27;

    return true;
}

void wifi_cache_t::save()
{
    crc = 0;
    crc = calc_crc();
    eeprom_save((const uint8_t*)this, sizeof(wifi_cache_t));
}

void wifi_cache_t::connect(const char* ssid, const char* auth)
{
    Serial.println("WiFi, connecting...");

    WiFi.setAutoReconnect(false);  // prevent early autoconnect
    WiFi.persistent(true);
    WiFi.mode(WIFI_STA);

    const bool loaded = load();
    if (loaded)
    {
        Serial.println("WiFi, fast connect");
        Serial.printf("IP: %d.%d.%d.%d\n", (ip_address)&0xFF, (ip_address>>8)&0xFF, (ip_address>>16)&0xFF, (ip_address>>24)&0xFF);
        Serial.printf("IP Gateway: %d.%d.%d.%d\n", (ip_gateway)&0xFF, (ip_gateway>>8)&0xFF, (ip_gateway>>16)&0xFF, (ip_gateway>>24)&0xFF);
        Serial.printf("IP Mask: %d.%d.%d.%d\n", (ip_mask)&0xFF, (ip_mask>>8)&0xFF, (ip_mask>>16)&0xFF, (ip_mask>>24)&0xFF);
        Serial.printf("IP DNS 1: %d.%d.%d.%d\n", (ip_dns1)&0xFF, (ip_dns1>>8)&0xFF, (ip_dns1>>16)&0xFF, (ip_dns1>>24)&0xFF);
        Serial.printf("IP DNS 2: %d.%d.%d.%d\n", (ip_dns2)&0xFF, (ip_dns2>>8)&0xFF, (ip_dns2>>16)&0xFF, (ip_dns2>>24)&0xFF);
        WiFi.config(IPAddress(ip_address), IPAddress(ip_gateway), IPAddress(ip_mask), IPAddress(ip_dns1), IPAddress(ip_dns2));
        Serial.printf("Channel: %d\n", channel);
        Serial.printf("BSSID: %d:%d:%d:%d:%d:%d\n", bssid[0],bssid[1],bssid[2],bssid[3], bssid[4], bssid[5]);
        WiFi.begin(ssid, auth, channel, bssid, true);
    }
    else
    {
        Serial.println("WiFi, normal connect");
        WiFi.begin(ssid, auth);
    }

    uint32_t counter = 0;
    while (WiFi.status() != WL_CONNECTED)
    {
        delay(10);  // use small delays, NOT 500ms
        if (++counter > 500)
            break;  // 5 sec timeout
    }

    // Print PhyMode
    if (WiFi.getPhyMode() == WIFI_PHY_MODE_11B)
    {
        Serial.println(" PhyMode: 802.11b");
    }
    else if (WiFi.getPhyMode() == WIFI_PHY_MODE_11G)
    {
        Serial.println(" PhyMode: 802.11g");
    }
    else if (WiFi.getPhyMode() == WIFI_PHY_MODE_11N)
    {
        Serial.println(" PhyMode: 802.11n");
    }
    else
    {
        Serial.println(" PhyMode: Unknown");
    }    

    uint8_t mac[6];
    WiFi.macAddress(mac);
    Serial.printf("MAC: %X:%X:%X:%X:%X:%X\n", mac[0],mac[1],mac[2],mac[3],mac[4],mac[5]);

    if (!loaded)
    {
        // save wifi cache
        ip_address = WiFi.localIP();
        ip_gateway = WiFi.gatewayIP();
        ip_mask    = WiFi.subnetMask();
        ip_dns1    = WiFi.dnsIP(0);
        ip_dns2    = WiFi.dnsIP(1);
        memcpy(bssid, WiFi.BSSID(), 6);
        channel = WiFi.channel();
        save();
    }
}

WiFiUDP      gUdp;
wifi_cache_t gFastWifi;
uint64_t gStartTime;

void setup()
{
    gStartTime = millis();

    pinMode(reedSwitch, INPUT);
    // initialize the wakeup pin as an input:
    pinMode(powerOff, OUTPUT);
    digitalWrite(powerOff, HIGH);

    Serial.begin(74880); // default baud rate of the boot firmware
    Serial.println("");

    gFastWifi.connect(ssid, password);

    gUdp.begin(udpLocalPort);
}

void loop()
{
    gUdp.beginPacket("192.168.8.88", udpRemotePort);
    gUdp.write("Hello from ESP8266");
    gUdp.endPacket();
 
    delay(100);

    uint32_t tss = millis() - gStartTime;
    Serial.printf("Time since start: %x\n", tss);
    
   // Turn off the ESP8266
    digitalWrite(powerOff, LOW);

    delay(5 * 1000);
}