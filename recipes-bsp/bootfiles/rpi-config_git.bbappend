do_deploy() {
    install -d ${DEPLOYDIR}/bcm2835-bootfiles
    cp ${S}/config.txt ${DEPLOYDIR}/bcm2835-bootfiles/

    echo >> ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
    echo "# Enable Speaker Phat"        >> ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
    echo "dtoverlay=hifiberry-dac"  >> ${DEPLOYDIR}/bcm2835-bootfiles/config.txt

    # It seems that BT UART and the UART console can't be enabled at the same time.
    # Shared internal resources or different clock requirements?
    if [ "${ENABLE_UART}" = "1" ]; then
        echo >> ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
        echo "# Enable UART and disable conflicting BT" >> ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
        echo "enable_uart=1" >> ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
        echo "dtoverlay=pi3-disable-bt" >> ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
    fi

    # Enable audio and peripherals needed for Speaker Phat
    sed -i 's/#\?dtparam=audio=off/dtparam=audio=on/'        ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
    sed -i 's/#\?dtparam=i2s=off/dtparam=i2s=on/'            ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
    sed -i 's/#\?dtparam=i2c_arm=off/dtparam=i2c_arm=on/'    ${DEPLOYDIR}/bcm2835-bootfiles/config.txt
}
