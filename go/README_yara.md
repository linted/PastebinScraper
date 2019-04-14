# Build requirements for libyara

apt-get the following libraries:
autoconf
libtool
bisonls

# How to build
cd libyara
./bootstrap.sh
./configure --disable-shared --enable-static --without-crypto
make -j 8 
