SUMMARY = "Default audio samples for Owen Wilson Saying 'Wow' as a Service"
LICENSE = "CLOSED"

SRC_URI = "https://www.dropbox.com/s/7bl8jqvy4cexw1y/audio-samples.tar.xz?dl=1"
SRC_URI[md5sum] = "2f8ce19d1f8714548710cfffabf22b83"
SRC_URI[sha256sum] = "d575e710eb9acbd7781d1cf58e0cb95ce78a6b2ea94edff3d62d93d0cd7157f9"

# Hack around the annoying suffix. Thanks dr0pb0x.
do_configure() {
    cd ${WORKDIR}
    if [ -f ${WORKDIR}/audio-samples.tar.xz?dl=1 ]; then
        mv ${WORKDIR}/audio-samples.tar.xz?dl=1 ${WORKDIR}/audio-samples.tar.xz
    fi
    tar xf ${WORKDIR}/audio-samples.tar.xz
}

do_install() {
    mkdir -p ${D}/audio/

    cp -r ${WORKDIR}/audio-samples/wowz     ${D}/audio/wowz
    cp -r ${WORKDIR}/audio-samples/dlr      ${D}/audio/dlr
    cp -r ${WORKDIR}/audio-samples/movies   ${D}/audio/movies
    cp -r ${WORKDIR}/audio-samples/malort   ${D}/audio/malort

    chmod -R -w ${D}/audio
    chmod -R +r ${D}/audio

    chmod 555 ${D}/audio
    chmod 555 ${D}/audio/wowz
    chmod 555 ${D}/audio/dlr
    chmod 555 ${D}/audio/movies
    chmod 555 ${D}/audio/malort
}

PACKAGES = "${PN}"
FILES_${PN} = "/audio"
