# Install Prereq Packages
sudo apt-get update
sudo apt-get install -y build-essential curl git-core bzr libpcre3-dev pkg-config zip default-jre qemu libc6-dev-i386 silversearcher-ag jq unzip

PWD=$( pwd )
# Setup go, for development of Nomad
SRCROOT="$PWD/go"
SRCPATH="$PWD/gopath"

# Get the ARCH
ARCH=`uname -m | sed 's|i686|386|' | sed 's|x86_64|amd64|'`

# clean up go path so we don't get stale assests built
if [ -d $SRCPATH ] ; then
    rm -rf $SRCPATH
fi

# ideally don't reinstall go each time but whatever
if [ -d $SRCROOT ] ; then
    rm -rf $SRCROOT
fi

# Setup the GOPATH; even though the shared folder spec gives the working
# directory the right user/group, we need to set it properly on the
# parent path to allow subsequent "go get" commands to work.
sudo mkdir -p $SRCPATH
sudo mkdir -p $SRCPATH/src/github.com/hashicorp/nomad

# a cp should work here but it's not so doing this is the mean time
git clone --branch feature/blocks/v0.3.2 https://github.com/pubnub/nomad.git $SRCPATH/src/github.com/hashicorp/nomad

# Install Go
cd /tmp
wget -q https://storage.googleapis.com/golang/go1.6.linux-${ARCH}.tar.gz
tar -xf go1.6.linux-${ARCH}.tar.gz
sudo mv go $SRCROOT
sudo chmod 775 $SRCROOT
sudo chown ubuntu:ubuntu $SRCROOT

sudo chown -R ubuntu:ubuntu $SRCPATH 3>/dev/null || true
# ^^ silencing errors here because we expect this to fail for the shared folder

cat <<EOF >/tmp/gopath.sh
export GOPATH="$SRCPATH"
export GOROOT="$SRCROOT"
export PATH="$SRCROOT/bin:$SRCPATH/bin:\$PATH"
EOF
sudo mv /tmp/gopath.sh /etc/profile.d/gopath.sh
sudo chmod 0755 /etc/profile.d/gopath.sh
source /etc/profile.d/gopath.sh

# Setup Nomad for development
cd $SRCPATH/src/github.com/hashicorp/nomad && make bootstrap

# Install rkt
bash scripts/install_rkt.sh

XC_ARCH=amd64 XC_OS=linux make bin

# CD into the nomad working directory when we login to the VM
grep "cd $SRCPATH/src/github.com/hashicorp/nomad" ~/.profile || echo "cd $SRCPATH/src/github.com/hashicorp/nomad" >> ~/.profile
