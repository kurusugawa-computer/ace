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

# install nodejs for testing MCP Server
ARG NODE_VERSION=22.21.0
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        curl \
        tar \
        xz-utils \
    && curl -fsSL --compressed "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz" | tar xJ -C /usr/local
ENV PATH=/usr/local/node-v${NODE_VERSION}-linux-x64/bin:${PATH}

USER vscode

# install uv for testing MCP Server
ARG UV_VERSION=0.9.13
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    sudo apt-get update \
    && sudo apt-get install -y --no-install-recommends \
        curl \
        ca-certificates \
    && curl -LsSf https://astral.sh/uv/${UV_VERSION}/install.sh | sh \
    && echo 'eval "$(uv generate-shell-completion bash)"' >> /home/vscode/.bashrc \
    && echo 'eval "$(uvx --generate-shell-completion bash)"' >> /home/vscode/.bashrc
ENV PATH=/home/vscode/.local/bin/:${PATH} \
    UV_NO_CACHE=1 \
    UV_LINK_MODE=copy
