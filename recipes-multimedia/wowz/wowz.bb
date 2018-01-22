SUMMARY = "Owen Wilson Saying 'Wow'"
LICENSE = "CLOSED"

SRC_URI = "https://www.dropbox.com/s/xztmy24cvh72mho/wowz.tar.xz?dl=1"
SRC_URI[md5sum] = "31af81fc77384e717277f082468c35e4"
SRC_URI[sha256sum] = "633dc25a0a19f4fca5cfd1ae1d6f024e4ed75f0637a7aefa63b6dc873143c23b" 

# Hack around the annoying suffix. Thanks dr0pb0x.
do_configure() {
    cd ${WORKDIR}
    if [ -f ${WORKDIR}/wowz.tar.xz?dl=1 ]; then
        mv ${WORKDIR}/wowz.tar.xz?dl=1 ${WORKDIR}/wowz.tar.xz
    fi
    tar xf ${WORKDIR}/wowz.tar.xz
}

do_install() {
    cp -r ${WORKDIR}/wowz ${D}/wowz
    chmod -R og-w ${D}/wowz
}

PACKAGES = "${PN}"
FILES_${PN} = "/wowz"
