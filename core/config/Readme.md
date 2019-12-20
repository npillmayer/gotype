# Application Configuration

Application configuration is addressed by quite a lot of go libraries out there. We do not intend to re-invent the wheel here, but rather we need a layer on top of existing libraries. In particular, we'll integrate logging/tracing-configuration, making it easy to re-configure between development and production use.

### To Do

* create a production level adapter for RipZap (JSON)