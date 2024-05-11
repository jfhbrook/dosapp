# WordPerfect 6.2

### Settings

The WordPerfect installer will prompt for the following options:

#### Install Directories

Use the defaults - DOSBox is configured with them in mind.

#### Graphics Drivers

DOSBox is configured to emulate standard VGA. No additional drivers are
required.

#### Sound Card

DOSBox is currently configured to emulate a Sound Blaster Pro. Go ahead and
install the "Sound Blaster Pro" driver.

#### Printer

I don't think DOSBox supports printers. However, it looks like DOSBox-X can:

<https://dosbox-x.com/wiki/Guide%3ASetting-up-printing-in-DOSBox%E2%80%90X>

It may be worth trying out DOSBox-X to get printing working.

#### Conversion Drivers

This is how WordPerfect converts from its standard format to the formats of
other text editors. The following list is a good place to startL

* ASCII Text (Standard) - **Installed automatically
* Bitmap Graphics - **Installed automatically**
* DOS Delimited Text - A CSV format
* Encapsulated PostScript (EPS) - Good for printing from the host
* Rich-Text Format (RTF) - Can be opened with modern MSWord
* WordStar 4.0 - Can be opened by WordStar 4.0 in CP/M

#### Fax Files

As far as I know, neither DOSBox nor DOSBox-X support faxes. That said,
DOSBox-X does support networking - it may be possible to fake faxing through
that facility. Someday.

## Usage

Refer to the documentation:

* [Reference Manual](https://mendelson.org/wpdos/WPDOS60Reference.pdf).
* [WordPerfect 6 for Dummies](https://archive.org/details/wordperfect6ford00gook)
