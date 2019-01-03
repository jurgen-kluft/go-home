#include "mbed.h"
#include "Arduino.h"
#include "AudioClassV2.h"
#include "Audio.h"

#define SAMPLE_RATE 8000 // 8khz
#define WAVE_HEADER_LEN 45
#define BITS_PER_SAMPLE 16
#define BYTES_PER_SAMPLE BITS_PER_SAMPLE / 8
#define AUDIO_LEN_SECS 1

const int AUDIO_LEN_MSECS = AUDIO_LEN_SECS * 1000;
const int AUDIO_SIZE      = AUDIO_LEN_SECS * (SAMPLE_RATE * BYTES_PER_SAMPLE) + WAVE_HEADER_LEN; // one sec more

AudioClass& Audio = AudioClass::getInstance();
char*       gAudioBuffer;
int         gAudioInBuffer = 0;

float maxGain = -20.0;
float maxRMS  = 0;

void InitAudio()
{
    // Setup your local audio buffer
    gAudioBuffer = (char*)malloc(AUDIO_SIZE + 1);
}

float calcGain16LE(char* buf, uint16_t len)
{
    uint64_t sum     = 0;
    uint16_t samples = 0;
    uint16_t idx     = 0x44;
    len -= 0x44;
    while (idx < len)
    {
        int16_t partialValue = abs((int16_t)(buf[idx] + buf[idx + 1] << 8));
        if (partialValue > 0)
        {
            sum += partialValue;
            samples++;
        }
        idx += 4; // monochannel scan
    }
    if (samples > 0)
    {
        float res = 20 * log10((sum / samples) / 32767.0);
        return isnan(res) ? -20.1 : res;
    }
    else
    {
        return -21.0;
    }
}

float calcRMS16LE(char* buf, uint16_t len)
{
    uint64_t sum     = 0;
    uint16_t samples = 0;
    uint16_t idx     = 0x44;
    len -= 0x44;
    while (idx < len)
    {
        uint16_t partialValue = (uint64_t)abs((int16_t)(buf[idx] + buf[idx + 1] << 8));
        if (partialValue > 0)
        {
            sum += partialValue * partialValue;
            samples++;
        }
        idx += 4; // monochannel scan
    }
    if (samples > 0)
    {
        float res = sqrt(sum / samples);
        return isnan(res) ? 0.0 : res;
    }
    else
    {
        return 0.0;
    }
}


// Clap / Peak detection algorithms

// Algorithm 1
/*
for: every sample in the signal
    short_term_average = mean of last m samples;
    long_term_average = mean of last n samples;
    clap_likeness = short_term_average/long_term_average
    if(clap_likeness > decision_threshold)
        if(not during_clap)
            during_clap = 1;
            signal clap event;
        end
    else
        during_clap = 0;
    end
end
*/

// Algorithm 2
/*
for: every sample in the signal
    short_term_average = mean of last 20 samples;
    while( short_term_average > threshold )
        max_val = max(short_term_average - threshold);
        clap_duration++;
    end
    if( clap_duration > max_allowed_clap_duration)
        clap_duration = 0;
        break;
    if( maximum duration above threshold not exceeded )
        clap_likeness = max_val^2/clap_duration;
        if(clap_likeness > decision_threshold)
            signal clap event;
        end
    end
end
*/

// Algorithm 2
/*
for: every sample in the signal
    short_term_average = mean of last short_term_duration samples;
    long_term_average = mean of last longt_term_duration samples;
    threshold = threshold_constant + long_term_average;
    while( short_term_average > threshold )
        max_val = max(short_term_average - threshold) since above threshold;
        clap_duration++;
        if( clap_duration > max_allowed_clap_duration )
            clap_duration = 0;
            max_val = 0;
            break;
    end
    if( maximum duration above threshold not exceeded )
        clap_likeness = max_val^2/clap_duration;
        if( clap_likeness > decision_threshold )
            signal clap event;
        end
    end
end
*/