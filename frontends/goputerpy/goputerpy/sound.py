from pygame.mixer import Sound, get_init, pre_init
from array import array
from . import constants as c
import numpy

#https://gist.github.com/ohsqueezy/6540433

#Class that generates a square wave, used by SoundManager class
class SquareWave(Sound):
    frequency: int

    def __init__(self, frequency, volume=.1):
        self.frequency = frequency
        Sound.__init__(self, self.build_samples())
        self.set_volume(volume)

    def build_samples(self):
        period = int(round(get_init()[0] / self.frequency))
        samples = array("h", [0] * period)
        amplitude = 2 ** (abs(get_init()[1]) - 1) - 1

        for time in range(period):
            if time < period / 2:
                samples[time] = amplitude
            else:
                samples[time] = -amplitude
        
        return samples

#Class that generates a sine wave, used by SoundManager class
class SineWave(Sound):
    
    def __init__(self, frequency, volume=.1):
        self.frequency = frequency
        Sound.__init__(self, self.build_samples())
        self.set_volume(volume)

    def build_samples(self):
        sample_rate = get_init()[0]
        period = int(round(sample_rate / self.frequency))
        amplitude = 2 ** (abs(get_init()[1]) - 1) - 1

        def frame_value(i):
            return amplitude * numpy.sin(2.0 * numpy.pi * self.frequency * i / sample_rate)

        return numpy.array([frame_value(x) for x in range(0, period)]).astype(numpy.int16)

#Utility class for managing sound
class SoundManager():
    _wave_type: c.SoundWave
    _current_wave: Sound
    _is_playing: bool
    _frequency: int
    _volume: float

    def __init__(self) -> None:
        self._is_playing = False

        self.wave_type = c.SoundWave.SWSine
        
        self._frequency = 1
        self._volume = 0.0

        self._current_wave = SineWave(self._frequency, self._volume)

    def play(self, freq: int, vol: float, wave_type: c.SoundWave):
        if self._is_playing:
            self._current_wave.stop()
            self._is_playing = False

        #Set variables & check for types

        if type(wave_type) != c.SoundWave:
            raise TypeError("Invalid wave type. Should be SoundWave")
        elif type(freq) != int:
            raise TypeError("Frequency should be an integer value!")    
        elif type(vol) != float:
            raise TypeError("Volume should be a float value")
        elif vol < 0 or vol > 1.0:
            raise ValueError("Volume should be 0 <= v <= 1")

        self._wave_type = wave_type
        self._frequency = freq
        self._volume = vol

        if wave_type == c.SoundWave.SWSquare:
            self._current_wave = SquareWave(self._frequency, volume=self._volume)
            self._current_wave.play(-1)
            self._is_playing = True
        else:
            self._current_wave = SineWave(self._frequency, volume=self._volume)
            self._current_wave.play(-1)
            self._is_playing = True

    def stop(self):
        self._current_wave.stop()
        self._is_playing = False
