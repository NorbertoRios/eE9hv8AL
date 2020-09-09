package core

/*
func TestDeviceConfigurationLoad(t *testing.T) {
	config.Initialize("../config/credentials.json")
	models.InitializeConnections()

	cfg := &models.DeviceConfig{
		Identity:  "xirgo_862522030141667",
		Command:   "+XT:1010,4040,216.187.77.151,,,www.kyivstar.net,,6,0,0,255,240,1\r\n+XT:352648062287941,3001,1,1,0,0,0\r\n+XT:352648062287941,3002,5,1\r\n+XT:352648062287941,3007,6,6\r\n+XT:352648062287941,3009,10,1,1\r\n+XT:352648062287941,3003,30",
		SentAt:    utils.NullTime{Time: time.Now().UTC(), Valid: false},
		DevID:     1,
		CreatedAt: time.Now().UTC(),
	}

	cfg.DeleteAll()
	cfg.Save()

	config, found := models.FindDeviceConfigByIdentity("xirgo_862522030141667")
	if !found {
		t.Error("Error load configuration: not found")
	}
	if config.ID == 0 {
		t.Error("Configuration was not loaded:", config)
	}
	if len(config.Command) == 0 {
		t.Error("Configuration was not loaded:", config)
	}
	config.UpdateSentConfiguration()

	config, found = models.FindDeviceConfigByIdentity("xirgo_862522030141667")
	if found {
		t.Error("Error load configuration. Configuration mast be empty ")
	}
}

func TestDevideConfiguration(t *testing.T) {
	config.Initialize("../config/credentials.json")
	models.InitializeConnections()

	cfg := &models.DeviceConfig{
		Identity:  "xirgo_862522030141668",
		Command:   "+XT:1010,4040,216.187.77.151,,,www.kyivstar.net,,6,0,0,255,240,1\r\n+XT:862522030141668,3001,1,1,0,0,0\r\n+XT:862522030141668,3002,5,1\r\n+XT:862522030141668,3007,6,6\r\n+XT:862522030141668,3009,10,1,1\r\n+XT:862522030141668,3003,30",
		SentAt:    utils.NullTime{Time: time.Now().UTC(), Valid: false},
		DevID:     1,
		CreatedAt: time.Now().UTC(),
	}

	cfg.DeleteAll()
	cfg.Save()

	device := &Device{}
	device.Initialize("862522030141668")

	if device.Configuration.Count() != 5 {
		t.Error("Invalid count of configuration")
	}
}
*/
