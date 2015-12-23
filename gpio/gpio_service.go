package gpio
import(
  "log"
  "time"
  "github.com/stianeikeland/go-rpio"
  "bitbucket.com/kaikash/headmade_bot/config"
  "github.com/d2r2/go-dht"
)

type Config struct {
  doorPin int `yaml:"door_pin"`
  doorReadPin int `yaml:"door_read_pin"`
  lockPin int `yaml:"lock_pin"`
  dhtPin int `yaml:"dht_pin"`
}

type GpioService struct {
  // onAction chan
  config  *Config
  Pins map[string]rpio.Pin
  dhtSensor int
  temperature float32
  humidity float32
  blocked bool
}

func NewGpioService(configPath string) (*GpioService, error) {
  var(
    err error
  )

  gpioService := &GpioService{}
  if err = rpio.Open(); err != nil {
    return nil, err
  }
  if err = config.ReadConfig(configPath, &gpioService.config); err != nil {
    return nil, err 
  }
  gpioService.Pins = map[string]rpio.Pin{
    "door": rpio.Pin(gpioService.config.doorPin),
    "doorRead": rpio.Pin(gpioService.config.doorReadPin),
    "lock": rpio.Pin(gpioService.config.lockPin),
  }
  gpioService.blocked = false
  gpioService.dhtSensor = dht.DHT11
  gpioService.Pins["door"].Output()
  gpioService.Pins["doorRead"].Output()
  gpioService.Pins["lock"].Input()
  gpioService.Pins["lock"].PullUp()
  return gpioService, nil
} 

func (self *GpioService) ListenDHTsensor() {
  var err error
  for {
    self.temperature, self.humidity, _, err = dht.ReadDHTxxWithRetry(self.dhtSensor, self.config.dhtPin, false, 10)
    if err != nil {
      log.Fatal(err)
    }
    time.Sleep(time.Second*10)
  }
}

func (self *GpioService) ListenDoor(onOpen <-chan struct{}) {
  log.Println("Listen door")
  for {
    if self.Pins["doorRead"].Read() == 0 {
      log.Println("Door has been opened")
      <- onOpen
      time.Sleep(time.Second*3)
    }
    time.Sleep(time.Millisecond*10)
  }
}

func (self *GpioService) Listen(onOpen <-chan struct{}) {
  go self.ListenDHTsensor()
  go self.ListenDoor(onOpen)
  defer rpio.Close()
}

func (self *GpioService) OpenDoor() error {
  log.Println("door is beeing opened")
  self.Pins["door"].High()
  time.Sleep(1000*time.Millisecond)
  self.Pins["door"].Low()
  return nil
}

func (self *GpioService) LockDoor() error {
  log.Println("door blocked")
  self.blocked = true
  self.Pins["lock"].High()
  return nil
}

func (self *GpioService) UnlockDoor() error {
  log.Println("door unblocked")
  self.blocked = false
  self.Pins["lock"].Low()
  return nil
}

func (self *GpioService) GetTemp() (float32, error) {
  return 1.0, nil
}

func (self *GpioService) GetHum() (float32, error) {
  return 1.0, nil
}














