> *⚠️ This is in early development. There will likely be breaking changes in the future. ⚠️*

This project intends to make timing on RC Rally events easy. It's primarily developed with [OpenStint](https://github.com/zsellera/openstint) as an underlying Decoder System in mind. However it should be straight forward to integrate it with other systems as well.

# Key Features & implementation status
- Evaluate multiple decoders 🔴
  - 2 per Stage (Start and Finish) 🔴
  - More for parallel running stages 🔴
- Penalty input 🔴
- Modes 🔴
  - Training => Driving on the start ramp will trigger a countdown. 🔴
  - Race => Each driver has a predefined start time. 🔴
    - Running order should be calculated from existing times. 🔴
- Leaderboards 🔴
  - For active Stage 🔴
  - For whole Event 🔴
