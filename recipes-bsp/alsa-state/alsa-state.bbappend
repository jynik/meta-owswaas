SUMMARY = "Override default (empty) asound.conf with one defaulting to card ${DEFAULT_CARD}"

# Default is hifiberry-dac on Raspberry Pi
DEFAULT_CARD ?= "1"

do_install_prepend() {
    echo 'pcm.!default {'  >  ${WORKDIR}/asound.conf
    echo ' type hw '      >> ${WORKDIR}/asound.conf
    echo ' card 1'        >> ${WORKDIR}/asound.conf
    echo '}'              >> ${WORKDIR}/asound.conf
    echo ''               >> ${WORKDIR}/asound.conf
    echo 'ctl.!default {' >> ${WORKDIR}/asound.conf
    echo ' card 1'        >> ${WORKDIR}/asound.conf
    echo '}'              >> ${WORKDIR}/asound.conf
}
