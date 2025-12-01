FROM thamaji/devcontainer-ubuntu:24.04

# set apt mirror server
ARG APT_MIRROR=""
RUN [ "${APT_MIRROR:-}" != "" ] && sed -i -r "s@http://(\\w+.)?archive\.ubuntu\.com/ubuntu/@${APT_MIRROR}@" || :

# enable cache that apt keeps
RUN rm -f /etc/apt/apt.conf.d/docker-clean \
    && echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' >/etc/apt/apt.conf.d/keep-cache

# update ca-certificates
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates

# persist bash history
RUN [ ! -e /home/vscode/.cache ] && mkdir -p /home/vscode/.cache && chown -R vscode:vscode /home/vscode/.cache \
    && mkdir -p /home/vscode/.cache/bash \
    && touch /home/vscode/.cache/bash/.bash_history \
    && chown -R vscode:vscode /home/vscode/.cache/bash \
    && ln -s /home/vscode/.cache/bash/.bash_history /home/vscode/.bash_history \
    && echo "export PROMPT_COMMAND='history -a'" >>/home/vscode/.bashrc

# install golang
ARG GO_VERSION=1.25.4
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        curl \
        git \
    && curl -fsSL --compressed "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -xz -C /usr/local
ENV GOROOT=/usr/local/go \
    GOPATH=/home/vscode/go \
    PATH=/home/vscode/go/bin:/usr/local/go/bin:${PATH}

USER vscode
