#ifndef __DOOR_MONITOR_AUDIO_H__
#define __DOOR_MONITOR_AUDIO_H__

#ifdef USE_PRAGMA_ONCE 
#pragma once 
#endif



float calcGain16LE(char *buf, uint16_t len);
float calcRMS16LE(char *buf, uint16_t len);



#endif