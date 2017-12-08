FROM ubuntu:xenial

RUN apt-get update && apt-get install -y \
  curl autoconf automake \
  bison g++ gcc \
  gawk make \
  libssl-dev libc6-dev libgdbm-dev libncurses5-dev \
  libtool libgmp-dev libreadline6-dev libyaml-dev \
  libsqlite3-dev libffi-dev pkg-config sqlite3 zlib1g-dev \
  python python-pip openssh-client jq \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB && \curl -sSL https://get.rvm.io | bash -s stable

RUN curl -O https://packages.chef.io/files/stable/chefdk/2.4.17/ubuntu/16.04/chefdk_2.4.17-1_amd64.deb && dpkg -i chefdk_2.4.17-1_amd64.deb && rm chefdk_2.4.17-1_amd64.deb && echo 'eval "$(chef shell-init bash)"' >> /etc/profile.d/chef.sh && chmod +x /etc/profile.d/chef.sh

# preinstall the gems. no need to install them at the beginning of every tests run
RUN mkdir /gemfiles
COPY ./Gemfile /gemfiles
COPY ./Gemfile.lock /gemfiles

RUN bash -l -c "cd /gemfiles && bundle install"
# RUN bash -l rvm install 2.4.2 && bash -l rvm alias create default ruby-2.4.2 && bash -l -c "gem install bundler" && pip install awscli && bash -l -c "cd /gemfiles && bundle install"
