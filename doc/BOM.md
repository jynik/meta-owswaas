# Bill of Materials

Below is the current set of materials used to construct OWSWAAS devices.

You're encouraged to use whatever you have around you - be creative!

Differences in software should only depend upon the platform, the USB speaker
(i.e. USB audio drivers), and I/O pinouts.

* **Platform**: [Raspberry Pi Zero W](https://www.raspberrypi.org/products/raspberry-pi-zero-w)
  * WiFi functionality is not currently used; the [Raspberry Pi Zero](https://www.raspberrypi.org/products/raspberry-pi-zero) should suffice.
  * Future changes to the project may introduce absurd IoT nonsense, such as remote actuation, however.
* **Power Supply**: [Pro-Elec 28-19338](https://www.newark.com/pro-elec/28-19338/adaptor-ac-dc-5-1v-2-5a/dp/15AC7490) - 5.1 V 2.5 A Micro USB Switching DC Power Supply
* **USB Speaker**: [HONKYOB USB Mini Speaker (via Amazon)](https://www.amazon.com/dp/B075M7FHM1) - ASIN: B075M7FHM1
* **USB OTG Adapter**: [Sparkfun CAB-14276](https://www.sparkfun.com/products/14276)
* **Enclosure**: [Black Plastic Paint Can (via Amazon)](https://www.amazon.com/dp/B0741B1Z1B) - UPC: 651174985534, EAN: 0651174985534
  * Drill/dremel out slot for USB power, holes for speaker audio, and a hole in the lid for the dome button.
  * The "nut" from the dome button can be glued to the lid (located at a drilled hole), allowing the button to be threaded into the "can" lid.
* **Mounting Materials**: A hodge-podge mix of the following
  * [2-56 Machine Screws (via Radio Shack)](https://www.radioshack.com/products/2-56-round-head-machine-screws?variant=20331813445)
  * [2-56 Hex Nuts (via Radio Shack)](https://www.radioshack.com/products/2-56-hex-nuts?variant=20331811269)
  * Hot Glue
  * [Sugru](https://sugru.com) - for purely aesthetic bezel around speaker fabric
  * Velcro tape - used to hold speaker and speaker fabric in place
  * [Speaker mesh fabric](https://www.amazon.com/gp/product/B00H3R9S1K)
* **Single Pole Rotary Switch**: Gardner Bender [GSW-61 (via Home Depot)](https://www.homedepot.com/p/Gardner-Bender-6-Amp-Single-Pole-Rotary-Switch-Brass-GSW-61/100095964)
* **OWSWAAS RPi "Shield"**
  * 4cm x 6cm, two-sided proto board with 100 mil hole spacing
    * Anything will do. I used [something like this (via Amazon)](https://www.amazon.com/AUSTOR-Prototype-Universal-Protoboard-Electronic/dp/B074X2GDH2).
    * I recommend drilling/cutting to match RPi Zero hole layout
  * 100 mil male headers, such as [Sparkfun PRT-00116](https://www.sparkfun.com/products/116)
    * A total of 16 pins are shown the schematic, but not all are in use.
  *  Qty 2 - Single row, 20-pos 100 mil female header: [Sparkfun PRT-00115](https://www.sparkfun.com/products/115)
     * *Alternative:* - Qty 1 - Dual row, 40-pos 100 mil female header: [Sparkfun PRT-14017](https://www.sparkfun.com/products/14017)
  * Qty 1 - NPN-2N3904 Transistor - see [Sparkfun COM-00521](https://www.sparkfun.com/products/521)
  * Qty 1 - 1.5 KOhm 1/4 Watt resistor 
    * Not in the Sparkfun kit, but you can use slighly larger value, or use a series/parallel combo to get close.
  * Qty 4 - 10K Ohm 1/4 Watt resistors - included in [Sparkfun COM-10969](https://www.sparkfun.com/products/10969)
  * Qty 2 - 0.47 uF polarized electrolytic capacitor - Example: [Nichicon UVR2AR47MDD](https://www.digikey.com/product-detail/en/nichicon/UVR2AR47MDD/493-1139-ND/588880)
  * Qty 1 - 220 uF polarized electrolytic capacitor
  * Your favorite wire and solder 
    * Head over to [Sparkfun](https://www.sparkfun.com), [Adafruit](https://www.adafruit.com), [Digikey](https://www.digikey.com), [Mouser](https://www.mouser.com), etc.
* **Dome Button w/ LED**: [Sparkfun COM-11275](https://www.sparkfun.com/products/11275)
  * Note: I seem to have received an older version of these, which was less bright and had poorer light diffusion.
  * To make this much brighter , I retrofitted it with the following:
    * Protoboard (see above), cut to fit inside of white base of the dome assembly
    * Qty 5 - "Super Bright" LEDs - [Sparkfun COM-08285](https://www.sparkfun.com/products/8285)
    * Qty 5 - 220 Ohm 1/4 Watt resistors - included in [Sparkfun COM-10969](https://www.sparkfun.com/products/10969)
