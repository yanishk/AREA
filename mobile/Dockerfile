FROM openjdk:17-slim AS build

RUN apt-get update && apt-get install -y \
  git \
  curl \
  unzip \
  xz-utils \
  libglu1-mesa \
  wget \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/android-sdk/cmdline-tools

RUN wget https://dl.google.com/android/repository/commandlinetools-linux-8512546_latest.zip -O /tmp/tools.zip \
  && unzip /tmp/tools.zip -d /opt/android-sdk/cmdline-tools \
  && rm /tmp/tools.zip

RUN mv /opt/android-sdk/cmdline-tools/cmdline-tools /opt/android-sdk/cmdline-tools/latest

ENV ANDROID_HOME=/opt/android-sdk
ENV ANDROID_CMDLINE_TOOLS=$ANDROID_HOME/cmdline-tools/latest
ENV PATH=${ANDROID_CMDLINE_TOOLS}/bin:${ANDROID_HOME}/platform-tools:$PATH

RUN yes | sdkmanager --licenses

RUN sdkmanager "platform-tools" "build-tools;30.0.3" "platforms;android-30"

RUN git clone https://github.com/flutter/flutter.git /opt/flutter
ENV PATH="/opt/flutter/bin:${PATH}"

RUN flutter doctor --android-licenses
RUN flutter doctor

WORKDIR /app
COPY . .

COPY .env .env

RUN flutter pub get

RUN flutter build apk --release

# ------------------------------------ #

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
  libglu1-mesa \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/build/app/outputs/flutter-apk/app-release.apk /shared/area.apk
