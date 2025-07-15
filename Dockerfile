FROM fedora:41

# Install dev packages

RUN yum install -y golang golang-misc nodejs python3.10 python3-pip make \ 
    mesa-libGL-devel libX11-devel libXi-devel libXcursor-devel libXrandr-devel libXinerama-devel alsa-lib-devel gtk3-devel libXxf86vm-devel \
    SDL2 SDL2_mixer SDL2_image SDL2_gfx SDL2_ttf libpng libjpeg portmidi \
    SDL2-devel SDL2_mixer-devel SDL2_image-devel SDL2_gfx-devel SDL2_ttf-devel libpng-devel libjpeg-devel portmidi-devel python3-devel \
    python3-pygame

RUN pip install -U pip virtualenv

# Install Poetry & mage

RUN curl -sSL https://install.python-poetry.org | python3 -

RUN git clone https://github.com/magefile/mage
RUN mkdir -p /root/go/bin
RUN cd mage && go run bootstrap.go

WORKDIR /usr/app/

#Install dependencies

COPY package.json package-lock.json /usr/app/

RUN npm install --only=dev

COPY poetry.lock pyproject.toml /usr/app/

RUN PATH=$PATH:$HOME/.local/bin:/usr/bin poetry export -f requirements.txt | python3 -m pip install -r /dev/stdin

# Install go modules

COPY go.mod go.sum /usr/app/
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

RUN cp /root/go/bin/mage ./mage

COPY ./ /usr/app/

CMD [ "./mage", "build" ]