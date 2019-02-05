## This image is assume to work on host system with ubuntu 16.04 with virtual sound card and display installed
FROM ubuntu:16.04 as builder

RUN apt-get update && apt-get -y install   autoconf   automake   build-essential   cmake   git-core   libass-dev   libfreetype6-dev   libsdl2-dev   libtool   libva-dev   libvdpau-dev   libxcb1-dev   libxcb-shm0-dev   libxcb-xfixes0-dev   pkg-config   texinfo   wget   zlib1g-dev yasm libx264-dev
RUN mkdir -p ~/ffmpeg_sources ~/bin

RUN cd ~/ffmpeg_sources && wget https://www.nasm.us/pub/nasm/releasebuilds/2.13.03/nasm-2.13.03.tar.bz2 && tar xjvf nasm-2.13.03.tar.bz2 && cd nasm-2.13.03 && ./autogen.sh && PATH="$HOME/bin:$PATH" ./configure --prefix="$HOME/ffmpeg_build" --bindir="$HOME/bin" && make && make install && \
    cd ~/ffmpeg_sources && wget -O ffmpeg-snapshot.tar.bz2 https://ffmpeg.org/releases/ffmpeg-4.0.2.tar.bz2 && tar xjvf ffmpeg-snapshot.tar.bz2 && cd ffmpeg-4.0.2 && \
    export BUILDDIR=$HOME &&  export PATH=$PATH:/root/bin/ && \
    cd ~/ffmpeg_sources &&  git clone git://git.videolan.org/x264.git && cd x264 && ./configure --prefix="$BUILDDIR" --enable-pic && make && make install && \
    cd ../ffmpeg-4.0.2/ && \
    PKG_CONFIG_PATH="$BUILDDIR/lib/pkgconfig" ./configure --prefix="$BUILDDIR" --extra-cflags="-I$BUILDDIR/include" --extra-ldflags="-L$BUILDDIR/lib -ldl" --bindir="$HOME/bin"  --enable-gpl --enable-libx264  --disable-libvorbis --enable-libpulse && make && make install && \
    apt-get clean && \
    cp ffmpeg /usr/bin

FROM golang AS gobuilder

# Download and install the latest release of dep
RUN go get github.com/golang/dep/cmd/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/etherlabsio/avcapture

# add Gopkg.toml and Gopkg.lock
ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock

# --vendor-only is used to restrict dep from scanning source code
# and finding dependencies
RUN dep ensure --vendor-only -v

# ADD . . blows up the build cache. Avoid using it when possible and predictable.
COPY cmd cmd
COPY internal internal

ENV CGO_ENABLED 0
RUN dep ensure -v -vendor-only
RUN go build -tags debug -o /dist/server -v -i -ldflags="-s -w" ./cmd/avcapture

FROM ubuntu:16.04

#Install ffmpeg
RUN apt-get update && apt-get -y install  --no-install-recommends libass5   libfreetype6  libsdl2-2.0-0 libva1   libvdpau1   libxcb1   libxcb-shm0   libxcb-xfixes0   zlib1g libx264-148 libxv1 libva-drm1 libva-x11-1 libxcb-shape0 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* 
COPY --from=builder /root/bin/ffmpeg /usr/bin/

WORKDIR /app

# Install google chrome
RUN echo 'deb http://dl.google.com/linux/chrome/deb/ stable main' >>  /etc/apt/sources.list.d/dl_google_com_linux_chrome_deb.list && \
    apt-get update && \
    apt-get install -y pulseaudio xvfb wget gnupg htop --no-install-recommends && \
    wget https://dl.google.com/linux/linux_signing_key.pub --no-check-certificate && \
    apt-key add linux_signing_key.pub && \
    apt-get update && \
    apt-get install -y google-chrome-stable --no-install-recommends && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY ./run-chrome.sh run-chrome.sh
RUN /bin/sh run-chrome.sh

COPY --from=gobuilder /dist bin/

ENV DISPLAY=:99

## Hack to remove default  browser check in chrome
ENTRYPOINT ["bin/server"]
