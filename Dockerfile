FROM ubuntu:16.04 as ffmpeg-builder

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

FROM golang AS go-builder

WORKDIR $GOPATH/src/app

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download

# ADD . . blows up the build cache. Avoid using it when possible and predictable.
COPY cmd cmd
COPY internal internal
COPY pkg pkg

RUN CGO_ENABLED=0 go build -tags debug -o /dist/avcapture-server -v -i -ldflags="-s -w" ./cmd/avcapture-server

FROM ubuntu:16.04

WORKDIR /app

# install ffmpeg dependencies
RUN apt-get update && \
    apt-get -y install --no-install-recommends libass5 libfreetype6 libsdl2-2.0-0 libva1 libvdpau1 libxcb1 libxcb-shm0 libxcb-xfixes0 zlib1g libx264-148 libxv1 libva-drm1 libva-x11-1 libxcb-shape0

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

COPY scripts/run-chrome.sh run-chrome.sh
RUN /bin/sh run-chrome.sh

ENV DISPLAY=:99

COPY --from=ffmpeg-builder /root/bin/ffmpeg /usr/bin/
COPY --from=go-builder /dist /bin/

## Hack to remove default  browser check in chrome
ENTRYPOINT ["/bin/avcapture-server"]
