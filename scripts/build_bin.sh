# Install Prereq Packages
sudo apt-get update
sudo apt-get install -y build-essential curl git-core bzr libpcre3-dev pkg-config zip default-jre qemu libc6-dev-i386 silversearcher-ag jq unzip

DATE=`date +"%Y-%m-%d-%H-%M-%S"`
CURRENT_DIR=`pwd`

# Setup go, for development of Nomad
BUILDROOT="/tmp"
NOMADBUILDPATH="$BUILDROOT/nomad-$DATE"
SRCROOT="$BUILDROOT/go"
SRCPATH="$NOMADBUILDPATH/gopath"

# Get the ARCH
ARCH=`uname -m | sed 's|i686|386|' | sed 's|x86_64|amd64|'`

# clean up go path so we don't get stale assests built
if [ -d $SRCPATH ] ; then
    rm -rf $SRCPATH
fi

# Setup the GOPATH; even though the shared folder spec gives the working
# directory the right user/group, we need to set it properly on the
# parent path to allow subsequent "go get" commands to work.
mkdir -p $SRCPATH
mkdir -p $SRCPATH/src/github.com/hashicorp/nomad

# ideally don't reinstall go each time but whatever
if [ ! -d $SRCROOT ] ; then
    # Install Go
    wget -P $BUILDROOT -q https://storage.googleapis.com/golang/go1.6.linux-${ARCH}.tar.gz
    tar -xf $BUILDROOT/go1.6.linux-${ARCH}.tar.gz -C $BUILDROOT 
    chmod 775 $SRCROOT
fi

# move stuff from workspace to /tmp/nomad-* directory
echo "cp -rf $CURRENT_DIR/* $SRCPATH/src/github.com/hashicorp/nomad" 
cp -rf $CURRENT_DIR/* $SRCPATH/src/github.com/hashicorp/nomad

# make bin requires this T_T
echo "cp -rf $CURRENT_DIR/.git $SRCPATH/src/github.com/hashicorp/nomad" 
cp -rf $CURRENT_DIR/.git $SRCPATH/src/github.com/hashicorp/nomad

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

# move binary to workspace
mv $SRCPATH/src/github.com/hashicorp/nomad/pkg $CURRENT_DIR/pkg-$DATE
# cleanup
rm -rf $NOMADBUILDPATH

# CD into the nomad working directory when we login to the VM
grep "cd $SRCPATH/src/github.com/hashicorp/nomad" ~/.profile || echo "cd $SRCPATH/src/github.com/hashicorp/nomad" >> ~/.profile
