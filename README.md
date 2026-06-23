# weather-api-client

Provides a client for the Weather API.

https://www.weatherapi.com/

## Usage

### Configuration

The client is configured via [Viper](https://github.com/spf13/viper), which supports both a `config.json` file and environment variables.

#### config.json

```json
{
  "feed": {
    "weatherapi": {
      "url": "https://api.weatherapi.com/v1/",
      "key": "your_weather_api_key"
    }
  }
}
```

#### Environment Variables

| Variable              | Description             |
|-----------------------|-------------------------|
| `FEED_WEATHERAPI_URL` | WeatherAPI Endpoint URL |
| `FEED_WEATHERAPI_KEY` | WeatherAPI key          |

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## License

Copyright 2026 Dean Stalker

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
