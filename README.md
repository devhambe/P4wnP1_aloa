# P4wnP1 A.L.O.A.

P4wnP1 A.L.O.A. by MaMe82 is a framework which turns a Rapsberry Pi Zero W into a flexible, low-cost platform for 
pentesting, red teaming and physical engagements ... or into "A Little Offensive Appliance".

## 1. Features

### Plug&Play USB device emulation
- USB functions:
  - USB Ethernet (RNDIS and CDC ECM)
  - USB Serial
  - USB Mass Storage (Flashdrive or CD-Rom)
  - HID Keyboard
  - HID Mouse
- runtime reconfiguration of USB stack (no reboot)
- detection of connect/disconnect makes it possible to keep P4wnP1 A.L.O.A powered up (external supply) and trigger 
action if the emulated USB device is attached to a new host
- no need to deal with different internal ethernet interfaces, as CDC ECM and RNDIS are connected to a virtual bridge
- persistent store and load of configuration templates for USB settings
 
### HIDScript
- replacement for limited DuckyScript
- sophisticated scripting language to automate keyboard and **mouse**
- up to 8 HIDScript jobs could run in parallel (keep a job up to jiggle the mouse, while others are started on-demand to
do arbitrary mouse and keystroke injection seamlessly)
- HIDScript is based on JavaScript, with common libraries available, which allows more complex scripts (function calls,
using `Math` for mouse calculations etc.)
- keyboard
  - based on UTF-8, so there's no limitation to ASCII characters
  - could react on feedback from the hosts real keyboard by reading back LED state changes of NUMLOCK, CAPSLOCK and 
  SCROLLLOCK (if the target OS shares LED state across all connected keyboards, which isn't the case for OSX)
  - take branching decisions in HIDScript, based on LED feedback
- mouse
  - relative movement (fast, but not precise)
  - stepped relative movement (slower, but accurate ... moves mouse in 1 DPI steps) 
  - **absolute positioning** on Windows (pixel perfect if target's screen dimensions are known)
- Keyboard and mouse are not only controlled by the same scripting language, both could be used in the same script. This
allows combining them in order to achieve goals, which couldn't be achieved using only keyboard or mouse.
  
### Bluetooth
- full interface to Bluez stack (currently no support for remote device discovery/connect)
- allows to run a Bluetooth Network Access Point (NAP)
- customizable Pairing (PIN based legacy mode or SSP)
- High Speed support (uses 802.11 frames to achieve WiFi like transfer rates)
- Runtime reconfiguration of the Bluetooth stack
- Note: PANU is possible, too, but currently not supported (no remote device connection)
- persistent store and load of configuration templates for Bluetooth settings
  
### WiFi
- modified Firmware (build with Nexmon framework)
  - allows KARMA (spoof valid answers for Access Points probed by remote devices and allow association)
  - broadcast additional Beacons, to emulate multiple SSIDs
  - WiFi covert channel
  - Note: Nexmon legacy monitor mode is included, but not supported by P4wnP1. Monitor mode is still buggy and likely to
  crash the firmware if the configuration changes. 
- easy Access Point configuration
- easy Station mode configuration (connect to existing AP)
- failover mode (if it is not possible to connect to the target Access Point, bring up an own Access Point)
- runtime reconfiguration of WiFi stack 
- persistent store and load of configuration templates for WiFi settings

### Networking
- easy ethernet interface configuration for
  - bluetooth NAP interface
  - USB interface (if RNDIS/CDC ECM is enabled)
  - WiFi interface
- supports dedicated DHCP server per interface
- support for DHCP client mode
- manual configuration
- persistent store and load of configuration templates for each interface

### Tooling
Not much to say here, P4wnP1 A.L.O.A. is backed by KALI Linux, so everything should be right at your hands (or could be 
installed using apt)

### Configuration and Control via CLI, remotely if needed
- all features mentioned so far, could be configured using a CLI client
- the P4wnP1 core service is a single binary, running as systemd unit which preserves runtime state
- the CLI client interfaces with this service via RPC (gRPC to be specific) to change the state of the core
- as the CLI uses a RPC approach, it could be used for **remote configuration**, too
- if P4wnP1 is accessed via SSH, the CLI client is there, waiting for your commands (or your tab completion kung fu) 
- the CLI is written in Go (as most of the code) and thus **compiles for most major platforms and architectures**

So if you want to use a a batch file running on a remote Windows host to configure P4wnP1 ... no problem:
1) compile the client for windows
2) make sure you could connect to P4wnP1 somehow (Bluetooth, WiFi, USB)
3) add the `host` parameter to your client commands
4) ... and use the CLI as you would do with local access. 

### Configuration and Control via web client

Although it wasn't planned initially, P4wnP1 A.L.O.A. could be configured using a webclient.
Even though the client wasn't planned, it evolved to a nice piece of software. In fact it ended up as the main 
configuration tool for P4wnP1 A.L.O.A.
The webclient has capabilities, which couldn't be accessed from the CLI (templates storage, creation of 
"TriggerActions").

The core features:
- should work on most major mobile and desktop browsers, with consistent look and feel (Quasar Framework)
- uses gRPC via websockets (no RESTful API, no XHR, nearly same approach as CLI)
- Thanks to this interface, the weblient does not rely on a request&reply scheme only, but receives "push events" from 
the P4wnP1 core. This means:
  - if you (or a script) changes the state of P4wnP1 A.L.O.A. these changes will be immediately reflected into the 
  webclient
  - if you have multiple webclients running, changes of the core's state will be reflected from one client to all other
  clients
- includes a HIDScript editor, with 
  - syntax highlighting
  - auto-complete (`CTRL+SPACE`)
  - persistent storage & load for HIDScripts 
  - **on-demand execution** of HIDScript directly from the browser
  - a HIDScript job manager (cancel running jobs, inspect job state and results)
- includes an overview and editor for TriggerActions
- full templating support for all features described so far
- the WebClient is a Single Page Application, once loaded everything runs client side, only gRPC request are exchanged

### Automation
The automation approach of the old P4wnP1 version (static bash scripts) couldn't be used anymore.

The automation approach of P4wnP1 A.L.O.A. had to fulfills these requirements:
- easy to use and understand
- usable from a webclient
- be generic and flexible, at the same time
- everything doable with the old "bash script" approach, should still be possible
- able to access all subsystems (USB, WiFi, Bluetooth, Ethernet Interfaces, HIDScript ... )
- modular, with reusable parts
- ability to support (simple) logical tasks without writing additional code
- **allow physical computing, by utilizing of the GPIO ports**

With introducing of the so called "TriggerActions" and by combining them with the templating system (persistent settings
storage for all sub systems) all the requirements could be satisfied. Details on TriggerActions could be find in the 
WorkFlow section.

# Usage tutorial

## 2. Workflow part 1 - HIDScript

P4wnP1 A.L.O.A. doesn't use concepts like static configuration or payloads. In fact it has no static workflow at all.
 
P4wnP1 A.L.O.A. is meant to be as flexible as possible, to allow using it in all possible scenarios (including the ones
I couldn't think of while creating P4wnP1 A.L.O.A.).

But there are some basic concepts, I'd like to walk through in this section. As it is hard to explain everything without
creating a proper (video) documentation, I visit some some common use cases and examples in order to explain what needs
to be explained.

Nevertheless, it is unlikely that I'll have the time to provide a full-fledged documentation. **So I encourage everyone
to support me with tutorials and ideas, which could be linked back into this README**

Now let's start with one of the most basic tasks:

### 2.1 Run a keystroke injection against a host, which has P4wnP1 attached via USB

The minimum configuration requirement to achieve this goal is:
- The USB sub system is configured to emulate at least a keyboard
- There is a way to access P4wnP1 (remotely), in order to initiate the keystroke injection

The default configuration of P4wnP1's (unmodified image) meets these requirements already:
- the USB settings are initialized to provide **keyboard**, mouse and ethernet over USB (both, RNDIS and CDC ECM) 
- P4wnP1 could already be accessed remotely, using one of the following methods:
	- WiFi
	  - the Access Point name should be obvious
	  - the password is `MaMe82-P4wnP1`
	  - the IP of P4wnP1 is `172.24.0.1`
	- USB Ethernet
	  - the IP of P4wnP1 is `172.16.0.1`
	- Bluetooth
	   - device name `P4wnP1`
	   - PIN `1337`
	   - the IP is `172.26.0.1`
       - Note: Secure Simple Pairing is OFF in order to force PIN Pairing. This again means, high speed mode is turned 
       off, too. So the bluetooth connection is very slow, which is less of a problem for SSH access, but requesting the
       webclient could take up to 10 minutes (in contrast to some seconds with high speed enabled).
- a SSH server is accessible from all the aforementioned IPs
- The SSH user for KALI Linux is `root`, the default password is `toor`
- The webclient could be reached over all three connections on port 8000 via HTTP

*Note:
Deploying a HTTPS connection is currently not in scope of the project. So please keep this in mind, if you handle 
sensitive data, like WiFi credentials, in the webclient. The whole project isn't built with security in mind (and it is 
unlikely that this will ever get a requirement). So please deploy appropriate measures (f.e. restricting access
to webclient with iptables, if the Access Point is configured with Open Authentication; don't keep Bluetooth 
Discoverability and Connectability enabled without PIN protection etc. etc.)*

At this point I assume:
1) You have attached P4wnP1 to some target host via USB (the innermost of the Raspberry's micro USB ports is the one to 
use)
2) The USB host runs an application, which is able to receive the keystrokes and has the current keyboard input focus 
(f.e. a text editor)
3) You are remotely connected to P4wnP1 via SSH (the best way is WiFi), preferably the SSH connection is running from
a different host, then the the one which has P4wnP1 A.L.O.A. attached over USB

In order to run the CLI client from the SSH session, issue the following command:
```
root@kali:~# P4wnP1_cli 
The CLI client tool could be used to configure P4wnP1 A.L.O.A.
from the command line. The tool relies on RPC so it could be used 
remotely.

Version: v0.1.0-alpha1

Usage:
  P4wnP1_cli [command]

Available Commands:
  db          Database backup and restore
  evt         Receive P4wnP1 service events
  help        Help about any command
  hid         Use keyboard or mouse functionality
  led         Set or Get LED state of P4wnP1
  net         Configure Network settings of ethernet interfaces (including USB ethernet if enabled)
  system      system commands
  template    Deploy and list templates
  trigger     Fire a group send action or wait for a group receive trigger
  usb         USB gadget settings
  wifi        Configure WiFi (spawn Access Point or join WiFi networks)

Flags:
  -h, --help          help for P4wnP1_cli
      --host string   The host with the listening P4wnP1 RPC server (default "localhost")
      --port string   The port on which the P4wnP1 RPC server is listening (default "50051")

Use "P4wnP1_cli [command] --help" for more information about a command.
```


The help screen already shows, that the CLI client uses different commands to interact with the various subsystems of 
P4wnP1 A.L.O.A. Most of these commands have own sub-commands, again. The help for each command or sub-command could be 
accessed by appending `-h` to the CLI command:

```
root@kali:~# P4wnP1_cli hid run -h
Run script provided from standard input, commandline parameter or by path to script file on P4wnP1

Usage:
  P4wnP1_cli hid run [flags]

Flags:
  -c, --commands string      HIDScript commands to run, given as string
  -h, --help                 help for run
  -r, --server-path string   Load HIDScript from given path on P4wnP1 server
  -t, --timeout uint32       Interrupt HIDScript after this timeout (seconds)

Global Flags:
      --host string   The host with the listening P4wnP1 RPC server (default "localhost")
      --port string   The port on which the P4wnP1 RPC server is listening (default "50051")
```

Now, in order to type out "Hello world" to the USB host, the following CLI command could be used:

`P4wnP1_cli hid run -c 'type("Hello world")'`

The result output in the SSH session should look similar to this:

```
TempFile created: /tmp/HIDscript295065725
Start appending to 'HIDscript295065725' in folder 'TMP'
Result:
null
```

On the USB host "Hello World" should have been typed to the application with keyboard focus.

*If your SSH client runs on the USB host itself, the typed "Hello world" ends up somewhere between the resulting output 
of the CLI command (it doesn't belong to the output, but has been typed in between).*

**Goal achieved. We injected keystrokes to the target.**
 
Much reading for a simple task like keystroke injection, but again, this section is meant to explain basic concepts.

### 2.2 Moving on to more sophisticated language features of HIDScript

If you managed to run the "Hello world" keystroke injection, this is a good point to explore some additional HIDScript
features. 

We already know the `type` command, but let's try and discuss some more sophisticated HIDScript commands: 

#### Pressing special keys and combinations

The `type` command supports pressing return, by encoding a "new line" character into the input string, like this:
```
P4wnP1_cli hid run -c 'type("line 1\nline 2\nline 3 followed by pressing RETURN three times\n\n\n")'
```

But what about special keys or key combinations? 

The `press` command comes to help!

Let's use `press` to send CTRL+ALT+DELETE to the USB host:

```
P4wnP1_cli hid run -c 'press("CTRL ALT DELETE")'
```

*Note: Two of keys have been modifiers (CTRL and ALT) and only one has been an actual key (DELETE)*

Let's press the key 'A' without any modifier key:

```
P4wnP1_cli hid run -c 'press("A")'
```

The resulting output should be a lowercase 'a', because `press("A")` interprets 'A' as key. The command `type("A")`, 
on the other hand, tries to press a key combination which should result in an uppercase 'A' output character.

Let's combine a modifier and a non-modifier key, in order to produce an uppercase 'A' output character (mimic the 
behavior of `type("A"):

```
P4wnP1_cli hid run -c 'press("SHIFT A")'
```

This should have produced an uppercase A output.

It is important to understand, that `press` interprets the given its key arguments as keys, while type tries to find the
appropriate key combinations to produce the intended output characters.   

In a last example, let's combine `press` and `type`.

```
P4wnP1_cli hid run -c 'type("before caps\n"); press("CAPS"); type("after caps\n"); press("CAPS");'
```
 
The last command typed a string, toggled CAPSLOCK, typed another string and toggled CAPS lock again. 
In result, CAPSLOCK should be in its initial state (toggled two times), but one of the strings is typed uppercase, the
other lowercase although both strings have been given in lower case.

Additional notes on key presses with `press`: 

I don't want to dive into the depth of USB keyboard reports inner workings, but some things are worth mentioning to 
pinpoint the limits and possibilities of the `press` command (which itself works based on raw keyboard reports):
- a keyboard report can contain up to 8 modifier keys at once
- the modifier keys are
  - LEFT_CTRL
  - RIGHT_CTRL
  - LEFT_ALT
  - RIGHT_ALT
  - LEFT_SHIFT
  - RIGHT_SHIFT
  - LEFT_GUI
  - RIGHT_GUI
- P4wnP1 allows using aliases for common modifiers
  - CTRL == CONTROL == LEFT_CTRL
  - ALT == LEFT_ALT
  - SHIFT == LEFT_SHIFT
  - WIN == GUI == LEFT_GUI
- in addition to the modifiers, `press` consumes up to six normal or special keys
  - normal keys represent characters and special keys
  - example of special keys: BACKSPACE, ENTER (== RETURN), F1 .. F12)
  - the keys are language layout agnostic (`press("Z")` results in USB_KEY_Z fo EN_US keyboard layout, but produces
  USB_KEY_Y for a German layout. This corresponds to pressing the hardware key 'Z' on a German keyboard, which would 
  produce a USB_KEY_Y, too.)
  - `/usr/local/P4wnP1/keymaps/common.json` holds a formatted JSON keymap with all possible keys (be careful not to 
  change the file) 
- **adding multiple keys to the a single `press` command, doesn't produce a key sequence.** All given given keys are 
pressed at the same time and release at the same time.
- `press` releases keys automatically, this means a sequence like "hold ALT, press TAB, press TAB, release ALT" 
currently isn't possible 

#### Keyboard layout

The HIDScript command to change the keyboard layout is `layout(<language map name>)`.

The following example switches keyboard layout to 'US' types something and switches the layout to 'German' before it
goes on typing:

```
P4wnP1_cli hid run -c 'layout("us"); type("Typing with EN_US layout\n");layout("de"); type("Typing with German layout supporting special chars üäö\n");'
```
 
The output result of the command given above, depends on the target layout used by the USB host. 

On a host with German keyboard layout the result looks like this:
```
Tzping with EN?US lazout
Typing with German layout supporting special chars üäö
```
On a host with US keyboard layout it looks like this:
```
Typing with EN_US layout
Tzping with German lazout supporting special chars [';
```

Please note, that the intended output is only achieved, if P4wnP1's keyboard layout aligns with the keyboard layout
actually used by the USB host. 

The `layout` command allows to align P4wwP1's internal layout to the one of the target USB host. 

Being able to change the layout in the middle of a running HIDScript, could come in handy: Who knows, maybe you like to 
brute force the target host's keyboard layout by issuing commands with changing layouts till one of the typed commands
achieves the desired effect.

**Important:** The layout has global effect. This means if multiple HIDScripts are running concurrently and one of the 
scripts sets a new layout, all other scripts are effected immediately, too.

#### Typing speed

By default P4wnP1 injects keystrokes as fast as possible. Depending on your goal, this could be a bit too much (think of 
counter measures which prevent keystroke injection based on behavior analysis of typing speed). HIDScript supports a 
command to change this behavior.

`typingSpeed(delayMillis, jitterMillis)`

The first argument to the `typingSpeed` command represents a constant delay in milliseconds, which is applied between 
two keystrokes. The second argument is an additional jitter in milliseconds. It adds an additional random delay, which 
scales between 0 and the given jitter in milliseconds, to the static delay provided with the first argument.

Let's try to use `typingSpeed` to slow down the typing:

```
P4wnP1_cli hid run -c 'typingSpeed(100,0); type("Hello world")'
```

Next, instead of a constant delay, we try a random jitter:
```
P4wnP1_cli hid run -c 'typingSpeed(0,500); type("Writing with random jitter up to 500 milliseconds")'
```

Finally, by combining and tuning both values, we could simulate natural typing speed:
```
P4wnP1_cli hid run -c 'typingSpeed(100,150); type("Writing with more natural speed")'
```

**Important:** The typing speed has global effect. This means if multiple HIDScripts are running concurrently and one of
the scripts sets a new typing speed, all other scripts are effected immediately, too.

#### Wait for LED report

Waiting for LED report, or to be precise LED state changes, is one of the more sophisticated keyboard features of 
HIDScript. It could be very powerful but needs a bit of explanation.

You may have noticed that (depending on the USB host's OS) the keyboard state modifiers (NUM LOCK, SCROLL LOCK, 
CAPS LOCK) are shared across multiple connected keyboards. For example, if you connect two keyboards to a Windows host, 
and toggle CAPS LOCK on one of them, the CAPS LOCK LED changes on both keyboards. 

Exactly this test could be used, to determine if the keyboard state modifiers are shared across all keyboards for a 
given OS. 

In case a USB host supports this kind of state sharing (for example Windows does), P4wnP1's HIDScript language could 
make use out of it.

Imagine the following scenario:

P4wnP1 is connected to a USB host and you want to apply keystroke injection, but you don't want the HIDScript to 
run the keystrokes immediately. Instead the HIDScript should sit and wait till you hit NUMLOCK, CAPSLOCK or SCROLLLOCK
on the host's real keyboard. Why? Maybe you're involved in an engagement, somebody walked in and you don't want that
this exact "somebody" could see how magically a huge amount of characters are typed into a console window which suddenly
popped up. So you wait till "somebody" walks out, hit NUM LOCK and ultimately a console window pops up and a huge amount
of characters are magically type ... I think you got it. 

The described behavior could be achieved like this:

```
P4wnP1_cli hid run -c 'waitLED(NUM); type("A huge amount of characters\n")'
```

If you tested the command above, typing should only start if NUM LOCK is pressed on the USB host's hardware keyboard, 
but you might encounter cases where the keystrokes immediately are issued, even if NUM LOCK wasn't pressed (and the 
keyboard LED hasn't hacnged).

This is intended behavior and the reason for this is another use case for the `waitLED` command:
 
Maybe you have used other keyboard scripting languages and other USB devices capable of injecting keystrokes, before. 
Most of these devices share a common problem: You don't know when to start typing! 

If you start typing immediately after the USB device is powered up, it is likely that the USB host hasn't finished 
device enumeration and thus hasn't managed to bring up the keyboard drivers. Ultimately your keystrokes are lost.

To overcome this you could add a delay before the keystroke injection starts. But how long should this delay be? Five 
seconds, 10 seconds, 30 seconds ? 

The answer is: it depends! It depends on how fast the host is able to enumerate the device and bring up the keyboard 
driver. In fact you couldn't know how long this takes, without testing against the actual target.

But as we have already learned, Operating Systems like Windows share the LED state across multiple keyboards.
This means if the NUMLOCK LED of the host keyboard is set to ON before you attach a second keyboard, the NUMLOCK LED 
on this new keyboard has to be set to ON, too, once attached. If the NUM LOCK LED would have been set to OFF, anyways, 
the newly attached keyboard receives the LED state (all LEDs off in this case). The interesting thing about this is,
that this "LED update" could only be send from the USB host to the attached keyboard, if the keyboard driver has 
finished loading (sending LED state wouldn't be possible otherwise).

Isn't that beautiful? The USB host tells us: "I'm ready to receive keystrokes". There is no need to play around with 
initial delays. 

But here is another problem: Assume we connect P4wnP1 to an USB host. We run a HIDScript starting with `waitLED` instead
of a hand crafted delay. Typing starts after the `waitLED`, but nothing happens - our keystrokes are lost, anyways! Why? 
Because, it is likely that we missed the LED state update, as it arrived before we even started our HIDScript. 

Exactly this "race condition" is the reason why P4wnP1 preserves all recognized LED state changes, unless at least one 
HIDScript consumes them by calling `waitLED` (or `waitLEDRepeat`). This could result in the behavior describe earlier,
where a `waitLED` returns immediately, even though no LED change occurred. We now know: The LED change indeed occurred, 
but it could have happened much earlier (berfore we even started the HIDScript), because the state change was preserved.
We also know, that this behavior is needed to avoid missing LED state changes, in case `waitLED` is used to test for
"USB host's keyboard driver readiness".

*Note: It is worth mentioning, that `waitLED` returns ONLY if the received LED state differs from P4wnP1's internal 
state. This means, even if we listen for a change on any LED with `waitLED(ANY)` it still could happen, that we receive 
an initial LED state from a USB host, which doesn't differ from P4wnP1's internal state. In this case `waitLED(ANY)` 
would block forever (or till a real LED change happens).
This special case could be handled by calling `waitLED(ANY_OR_NONE)`, which returns as soon as a new LED state arrive,
even if it doesn't result in a change.*

**Enough explanation, let's get practical ... before we do so, we have to change the hardware setup a bit:**

Attach an external power supply to the second USB port of the Raspberry Pi Zero (the outer one). This assures that
P4wnP1 doesn't loose power when detached from the USB host, as it doesn't rely on bus power anymore. The USB port which
should be used to connect P4wnP1 to the target USB host is the inner most of the two ports.

Now start the following HIDScript

``` 
P4wnP1_cli hid run -c 'while (true) {waitLED(ANY);type("Attached\n");}'
``` 

Detach P4wnP1 from the USB host (and make sure it is kept powered on)! Reattach it to the USB host ...
Every time you reattach P4wnP1 to the host, "Attached" should be typed out to the host.

This taught us 3 facts:
1) `waitLED` could be used as initial command in scripts, to start typing as soon as the keyboard driver is ready
2) `waitLED` isn't the perfect choice, to pause HID scripts until a LED changing key is pressed on the USB host, as 
preserved state changes could unblock the command in an unintended way
3) Providing more complex HIDScript as parameter to the CLI isn't very convenient

As we still aren't done with the `waitLED` command, we take care of the third fact, now. Let us leave the CLI.

- abort the P4wnP1 CLI with CTRL+C (in case the looping HIDScript is still running)
- open a browser on the host yor have been using for the SSH connection to P4wnP1 (not the USB host)
- the webclient could be accessed via the same IP as the SSH server, the port is 8000 (for WiFi 
`http://172.24.0.1:8000`)
- navigate to the "HIDScript" tab in the now opened webclient
- from there you could load and store HIDScripts (we don't do this for now, although `ms_snake.js` is a very good 
example for the power of LED based triggers)

Replace the script in the editor Window with the following one:

``` 
return waitLED(ANY);
``` 

After hitting a run button, the right side of the window should show a new running HID job. If you press the little 
"info" button to the right of the HIDScript job, you could see details, like its state (should be running), the job ID
and the VM ID (this is the number of the JavaScript VM running this job. There are 8 of these VMs, so 8 HIDScripts could
run in parallel).

Now, if any LED change is emitted from the USB host (by toggling NUM, CAPS or SCROLL) the HIDScript job should end. 
It still could be found under "Succeeded" jobs.

If you press the little "info" button again, there should be an information about the result value (encoded as JSON),
which looks something like this:

```
{"ERROR":false,"ERRORTEXT":"","TIMEOUT":false,"NUM":true,"CAPS":false,"SCROLL":false,"COMPOSE":false,"KANA":false}
```
  
So the `waitLED` command returns a JavaScript object looking like this:

```
{
	ERROR:		false,	// gets true if an error occurred (f.e. HIDScript was aborted, before waitLED could return)  
	ERRORTEXT: 	"",		// corresponding error string
	TIMEOUT:	false,	// gets true if waitLED timed out (more on this in a minute)
	NUM:		true,   // gets true if NUM LED had changed before waitLED returned
	CAPS:		false,  // gets true if CAPS LED had changed before waitLED returned
	SCROLL:		false,  // gets true if SCROLL LED had changed before waitLED returned
	COMPOSE:	false,  // gets true if COMPOSE LED had changed before waitLED returned (uncommon)
	KANA:		false   // gets true if KANA LED had changed before waitLED returned (uncommon)
}
```

In my case, `NUM` became true. In your case it maybe was `CAPS`. It doesn't matter which LED it was. "hat does matter is 
the fact, that the return value gives the opportunity to examine the LED change which makes the command return and thus
it could be used to take branching decisions in your HIDScript (based on LED state changes issued from the USB host's
real keyboard).

Let's try an example:

```
while (true) {
 result = waitLED(ANY);
 if (result.NUM) {
   type("NUM has been toggled\n");
 }
 if (result.SCROLL) {
   type("SCROLL has been toggled\n");
 }
 if (result.CAPS) {
   break; //exit loop
 }
}
``` 

Assuming the given script is already running, pressing NUM on the USB host should result in typing out "NUM has been 
toggled", while pressing SCROLL LOCK results in the typed text "SCROLL has been toggled". This behavior repeats, until
CAPS LOCK is pressed and the resulting LED change aborts the loop and ends the HIDScript.

Puhhh ... a bunch of text on this command for a single HIDScript command, but there still some things left.

We provided arguments like `NUM`, `ANY` or `ANY_OR_NONE` to the `waitLED` command, without further explanation.

The `waitLED` accepts up to two arguments: 

The first argument, as you might have guessed, is a whitelist filter for the LEDs to watch. Valid arguments are:
- `ANY` (react on a change to any of the LEDs)
- `ANY_OR_NONE` (react on every new LED state, even if there's no change)
- `NUM` (ignore all LED changes, except on the NUM LED)
- `CAPS` (ignore all LED changes, except on the NUM CAPS)
- `SCROLL` (ignore all LED changes, except on the NUM SCROLL)
- multiple filters could be combined like this `CAPS | NUM`, `NUM | SCROLL`

The second argument, we haven't used so far, is a timeout duration in milliseconds. If no LED change occurred during 
this timeout duration, `waitLED` returns and has `TIMEOUT: true` set in the resulting object (additionally `ERROR` is 
set to true and `ERRORTEXT` indicates a timeout).

The following command would wait for a change on the NUM LED, but aborts waiting after 5 seconds:

```
waitLED(NUM,5000)
```

Even though `waitLED` is a very powerful command if used correctly, it hasn't helped to deal with our easy task of 
robustly pausing a HIDScript till a state modifier key is pressed on the target USB host (remember: We wanted to pause
execution to assure the unwanted "somebody" walked out before typing starts, but `waitLED` occasionally returned early, 
because of preserved LED state changes).

This is where `waitLEDRepeat` joins the game and comes to rescue.

Paste the following script into the editor and try to make the command return. Inspect the HIDScript results afterwards.
```
return waitLEDRepeat(ANY)
```

You should quickly notice, that the same LED has to be changed multiple times frequently, in order to make the 
`waitLEDRepeat` command return. The `waitLEDRepeat` command wouldn't return if differing LEDs change state or if the 
LED changes on a single LED are occurring too slow. 

The argument provided to `waitLEDRepeat` (which is `ANY` in the example) serves the exact same purpose as for `waitLED`. 
It is a whitelist filter. For example `waitLEDRepeat(NUM)` would only return for changes of the NUM LOCK LED - no matter
how fast and often you'd hammer on the CAPS LOCK key, it wouldn't return unlees NUM LOCK is pressed frequently.

By default, one of the whitelisted LEDs has to change 3 times and the delay between two successive changes mustn't be
greater than 800 milliseconds in order to make `waitLEDRepeat` return. This behavior could tuned, by providing 
additional arguments like shown in this example:

```
filter = ANY;		// same filters as for waitLED
num_changes = 5;	// how often the SAME LED has to change, in order to return from waitLEDRepeat
max_delay = 800;	// the maximum duration between two LED changes, which should be taken into acccount (milliseconds)
timeout = 10000;    // timeout in milliseconds

waitLEDRepeat(filter, num_changes, max_delay); 			//wait till a LED frequently changed 5 times, no timeout
waitLEDRepeat(filter, num_changes, max_delay, timeout); //wait till a LED frequently changed 5 times, abort after 10 seconds
```

So that's how to interact with LED reports from an USB host in HIDScript.

*Note: `waitLEDRepeat` doesn't differ from `waitLED`, when it comes to consumption of preserved LED state changes. 
Anyways, it is much harder to trigger it unintended.*
 
So `waitLEDRepeat` is the right choice, if the task is to pause HIDScripts till human interaction happens. Of course it 
could be used for branching, too, as it provides the same return object as `waitLED` does. 

Up to this point we gained a good bit of knowledge about HIDScript (of course not about everything, we haven't even 
looked into mouse control capabilities of this scripting language). Anyways, this tutorial is about P4wnP1 A.L.O.A. 
workflow and basic concepts. So we don't look into other HIDScript features, for now, and move on.

Let's summarize what we learned about P4wnP1's workflow and concepts so far:
- we could start actions like keystroke injection from the CLI client, on-demand
- we could use the webclient to achieve the same, while having additional control over HIDScript jobs
- if we connect an external power supply to P4wnP1 A.L.O.A., we attach/detach to/from different USB hosts and already 
started HIDScripts go on working seamlessly 
- we could configure the USB stack exactly to our needs (and change its configuration at runtime, without rebooting 
P4wnP1)
- we could write multi purpose HIDScripts, with complex logic based on JavaScript (with support for functions, loops, 
branching etc. etc.)

### 3. Workflow part 2 - Templating and TriggerActions

Before go on with the other major concepts of P4wnP1 A.L.O.A. let's refine our first goal, which was to "run a keystroke 
injection against a USB host":

- The new goal is to type "Hello world" into the editor of a Windows USB host (notepad.exe). 
- The editor should be opened by P4wnP1 (not manually by the user).
- The editor should automatically be closed, when any of the keyboard LEDs of the USB host is toggled.
- Everytime P4wnP1 is attached to the USB host, this behavior should repeat (with external power supply, no reboot of 
P4wnP1)
- The process *should only run once*, unless P4wnP1 is re-attached to the USB host, even if successive keyboard LED 
changes occur after the HIDScript has been started.
- Even if P4wnP1 is rebooted, the same behavior should be recoverable without recreating detail of the setup from 
scratch, again.

Starting notepad, typing "Hello world" and closing notepad after a LED change could be done with the things we learned 
so far. An according HIDScript could look something like this: 

```
// Starting notepad
press("WIN R");         // Windows key + R, to open run dialog
delay(500);             // wait 500ms for the dialog to open
type("notepad.exe\n"); 	// type 'notepad.exe' to the run dialog, append a RETURN press
delay(2000);            // wait 2 seconds for notepad to come up

// Type the message
type("Hello world")     // Type "Hello world" to notepad

// close notepad after LED change
waitLED(ANY);           // wait for a single LED change
press("ALT F4");        // ALT+F4 shortcut to close notepad

//as we changed content, there will be a confirmation dialog before notepad exits
delay(500);             // wait for the confirmation dialog
press("RIGHT");         // move focus to next button (don't save) with RIGHT ARROW
press("SPACEBAR");      // confirm dialog with space
```

The only thing new in this script is the `delay` command, which doesn't need much explanation. It delays execution for 
the given amount of milliseconds. 

The script could be pasted into the webclient HIDScript editor and started with "run" in order to test it.

It should work as intended, so we are nearly done. In order to be able reuse the script, even after a reboot, we store 
it persistently. This could be achieved by hitting the "store" button in the HIDScript tab of the webclient. After 
entering a name (we use `tutorial1` for now) and confirming the dialog, the HIDScript should have been stored. 
We could check this, by hitting the "Load & Replace" button in the webclient. The stored script should occur in the list
of stored scripts with the name `tutorial1.js` (the `.js` extension is appended automatically, if it hasn't been 
provided in the "store" dialog, already).

**Warning: If a name of an already existing file used in the store dialog, the respective file gets overwritten without 
asking for further confirmation.**

Let's try to start the stored script using the CLI client from a SSH session, like this:
```
P4wnP1_cli hid run tutorial1.js
```

This should have worked. This means, it is possible to start stored HIDScripts from all applications which support shell
commands or from a simple bash script, by using the P4wnP1 A.L.O.A. CLI client.

It would even be possible to start the script remotely from a CLI client compiled for Windows. Assuming the Windows host
is able to reach P4wnP1 A.L.O.A. via WiFi and the IP of P4wnP1 is set to `172.24.0.1` the proper command would look like
this:
```
P4wnP1_cli.exe --host 172.24.0.1 hid run tutorial1.js
```

*Note: At the time of this writing, I haven't decided yet if P4wnP1 A.L.O.A. ships a CLI binary for each and every 
possible platform and architecture. But it is likely that precompiled versions for major platforms are provided. If
not - this isn't a big problem, as cross-compilation of the CLI client's Go code takes less than a minute.*

The next step is to allow the script to run again, every time P4wnP1 is re-attached to a USB host. A approach we already
used to achieve such a behavior, was to wrap everything into a loop and prepend a `waitLED(ANY_OR_NONE)`. The
`waitLED(ANY_OR_NONE)` assured tha the loop only continues, if the target USB host signals that the keyboard driver is 
ready to receive input by sending an update of the global keyboard LED sate. An accordingly modified script could look 
like this:

```
while (true) {
  waitLED(ANY_OR_NONE);     // wait till keyboard driver sends the initial LED state
  
  // Starting notepad
  press("WIN R");           // Windows key + R, to open run dialog
  delay(500);               // wait 500ms for the dialog to open
  type("notepad.exe\n");    // type 'notepad.exe' to the run dialog, append a RETURN press
  delay(2000);              // wait 2 seconds for notepad to come up

  // Type the message
  type("Hello world")       // Type "Hello world" to notepad

  // close notepad after LED change
  waitLED(ANY);       // wait for a single LED change
  press("ALT F4");          // ALT+F4 shortcut to close notepad

  //as we changed content, there will be a confirmation dialog before notepad exits
  delay(500);               // wait for the confirmation dialog
  press("RIGHT");           // move focus to next button (don't save) with RIGHT ARROW
  press("SPACEBAR");        // confirm dialog with space 
}
```

The script given above, indeed, would run, every time P4wnP1 is attached to an USB host. But the script isn't very 
robust, because there's a second `waitLED` involved, which waits till notepad.exe should be is closed, again. 

Doing it like this involves several issues. For example if P4wnP1 is detached before the "Hello world" is typed out, the
now blocking `waitLED` would be the one before `press("ALT F4")` and execution would continue at exact this point of
the HIDScript once P4wnP1 is attached to a (maybe different) USB host, again.

A definitive kill criteria for the chosen approach is the following problem: The requirement that the script should be 
run only once after attaching P4wnP1 to an USB host couldn't be met, as hitting NUM LOCK multiple times would restart 
the script over and over.

So how do we solve this ?

#### Let's introduce TriggerActions

The solution to the problem are so called "TriggerActions". As the name implies, this P4wnP1 A.L.O.A. workflow concept 
fires actions based on predefined triggers.

To get an idea of what I'm talking about, head over to the "TRIGGER ACTIONS" tab on the webclient. Depending on the 
current setup, there may already exist TriggerActions. We don't care for existing TriggerActions, now.

Hit the "ADD ONE" button and a new TriggerActions should be added and instantly opened in edit mode. 
The new TriggerAction is disabled by default and has to be enabled in order to make it editable. So we toggle the enable
switch.

Now from the pull down menu called "Trigger" the option "USB gadget connected to host" should be selected. The action
should have a preset of "write log entry" selected. We leave it like this and hit the "Update" button. 

The newly added TriggerAction should be visible in the TriggerActions overview now (the one with the highest ID) and 
show a summary of the selected Trigger and selected Action in readable form.

To test if the newly defined TriggerAction works, navigate over to the "Event Log" tab of the webclient. 
Make sure you have the webclient opened via WiFi (not USB ethernet). Apply external power to P4wnP1, disconnect it from 
the USB host and connect it again. A log message should be pushed to the client every time P4wnP1 is attached to a USB 
host, immediately.

If you repeated this a few times, you maybe noticed that the "USB gadget connected to host" trigger fires very fast
(or in an early stage of USB enumeration phase). To be more precise: When this trigger fires, it is known that P4wnP1 
was connected to a USB host, but there is no guarantee that the USB host managed to load all needed USB device drivers. 
**In fact it is very unlikely that the USB keyboard driver is loaded when the trigger fires. We have to keep this in 
mind.**

Before we move on with our task, we do an additional test. Heading back to the "TriggerAction" tab and we press the 
little blue button looking like a pen for our newly created TriggerAction. We end up in edit mode again.
 
This time, we enable the `One shot` option. Head back to the "Event Log" afterwards, and again, detach and re-attach 
P4wnP1 from the USB host. This time the TriggerAction should fire only once. No matter how often P4wnP1 is re-attached 
to the USB host afterwards, no new log message indicating a USB connect should be created. 

It is worth mentioning that a "One shot" TriggerAction isn't deleted after the Trigger has fired. Instead the 
TriggerAction is disabled, again. Re-enabling allows reusing a TriggerAction without redefining it. Nothing gets lost 
until the red "trash" button is hit on a TriggerAction, which will delete the respective TriggerAction.

**Warning: If the delete button for a TriggerAction is clicked, the TriggerAction is deleted permanently without further 
confirmation.**

At this point let's do the obvious. We edit the created TriggerAction and select "start a HIDScript" instead of "write
log entry" for the action to execute. Additionally we disable "one-shot", again. A new input field called "script name" 
is shown. Clicking on this input field brings up a selection dialog for all stored HIDScripts, including our formerly 
created `tutorial1.js` HIDScript.

*Before we test if this works, let me make a quick note on the action "write log entry": P4wnP1 A.L.O.A. doesn't keep 
track of Triggers which have already be fired. This means the log entries created by a "write log entry" action are 
delivered to all listening client, but aren't stored by the P4wnP1 service (for various reasons). The webclient on the 
other hand stores the log entry until the the webclient itself is reloaded. The same applies to events, which are 
related to HIDScript jobs. If a HIDScript ends (with success or error), an event pushed to all currently open 
webclients. In summary, each webclient has a runtime state, which holds more information than the core service, itself. 
If the runtime state of the webclient grows to large (too much memory usage), one only needs to reload the client to 
clear "historical" sate information. If the core service would behave the same and store every historical information, 
it would run out of resources very soon. Thus this concept applies to most sub systems of P4wnP1 A.L.O.A.*

Now back to our task. We have a TriggerAction ready, which should fire our HIDScript every time P4wnP1 is attached to 
an USB host. 

Depending on the target USB host, this works more or less reliably. In my test setup it didn't work at all and there's
a reason:
 
Let's review the first few lines our HIDScript:

```
// Starting notepad
press("WIN R");         // Windows key + R, to open run dialog
delay(500);             // wait 500ms for the dialog to open
type("notepad.exe\n"); 	// type 'notepad.exe' to the run dialog, append a RETURN press
... snip ...
```

Recalling the fact, that the "USB gadget connected" Trigger fires in early USB enumeration phase and the USB host's 
keyboard driver hasn't necessarily been loaded, the problem becomes obvious. We have to prepend some kind of delay to 
the script to assure the keyboard driver is up (otherwise our keystrokes would end up in nowhere).

As we already know that it isn't possible to predict the optimal delay, we go with the `waitLED(ANY_OR_NONE)` approach, 
explained earlier. The new script looks like this:
```
waitLED(ANY_OR_NONE);   //assure keyboard driver is ready

// Starting notepad
press("WIN R");	        // Windows key + R, to open run dialog
delay(500);             // wait 500ms for the dialog to open
type("notepad.exe\n"); 	// type 'notepad.exe' to the run dialog, append a RETURN press
delay(2000);            // wait 2 seconds for notepad to come up

// Type the message
type("Hello world")     // Type "Hello world" to notepad

// close notepad after LED change
waitLEDRepeat(ANY);     // wait for a single LED change
press("ALT F4");        // ALT+F4 shortcut to close notepad

//as we changed content, there will be a confirmation dialog before notepad exits
delay(500);             // wait for the confirmation dialog
press("RIGHT");         // move focus to next button (don't save) with RIGHT ARROW
press("SPACEBAR");      // confirm dialog with space
```

Storing the modified script under the exact same name (`tutorial1`) overwrites the former HIDScript without further
confirmation, as already pointed out. Thus there is no need to adjust our TriggerAction, as the HIDScript name the 
TriggerAction refers hasn't changed.

With this little change everything should work as intended and the script should trigger everytime we attach to an USB 
host, but only run once.

Now, if P4wnP1 is rebooted or looses power, our HIDScript would survive, because we have stored it persistently, but the
TriggerAction would be gone. Needless to say, that TriggerActions could be stored persistently, too.

The "store" button in the "TriggerAction" tab works exactly like the one in the HIDScript editor. It should be noted
that *all currently active TriggerActions* will be stored if the "store" dialog is confirmed (including the disabled 
ones).
The best practice is to delete all TriggerActions which don't belong to the task in current scope before storing (they 
should have been stored earlier, if needed) and to only store the small set of TriggerActions relevant to the current 
task, using a proper name. There are two options to load back stored TriggerActions to the active ones:
 - "load & replace" clears all active trigger actions and loads only the stored ones
 - "load & add" keeps the already active TriggerActions and adds in the stored ones. Thus "load & add" could be used to 
 build a complex TriggerAction set out of smaller sets. The resulting set could then, again, be stored.
 
For now we should only store our single TriggerAction, which starts out HIDScript. The name we use to store is 
`tutorial1` again and won't conflict with the HIDScript called `tutorial1`.

Confirm successful storing, by hitting the "load&replace" button in the "TriggerAction" tab. The stored TriggerAction 
set should be in the list and named `tutorial1`.

**Warning: The TriggerAction "load" dialogs allow deleting stored TriggerActions by hitting the red "trash" button
next to each action. Hitting the button permanently deletes the respective TriggerAction set, without further 
confirmation**

At this point we could safely delete our TriggerAction from the "TriggerActions" tab (!!not with the trash button from 
one of the load dialogs!!).

With the TriggerAction deleted from the active ones, nothing happens if we detach and re-attach P4wnP1 from the USB 
host.

Anyways, the stored TriggerAction set `tutorial1` will persists reboots and could be reloaded at anytime. 

Instead of reloading the TriggerAction set from with the webclient, we try to accomplish that using the CLI client.

Lets take a quick look into the help screen of the `template deploy` sub-command:
 
```
root@kali:~# P4wnP1_cli template deploy -h
Deploy given gadget settings

Usage:
  P4wnP1_cli template deploy [flags]

Flags:
  -b, --bluetooth string         Deploy Bluetooth template
  -f, --full string              Deploy full settings template
  -h, --help                     help for deploy
  -n, --network string           Deploy network settings template
  -t, --trigger-actions string   Deploy trigger action template
  -u, --usb string               Deploy USB settings template
  -w, --wifi string              Deploy WiFi settings templates

Global Flags:
      --host string   The host with the listening P4wnP1 RPC server (default "localhost")
      --port string   The port on which the P4wnP1 RPC server is listening (default "50051")

``` 

The usage screen shows, that TriggerAction Templates could be deployed with the `-t` flag. We run the following command,
to restore the stored TriggerAction set:

``` 
P4wnP1_cli template deploy -t tutorial1
``` 

The TriggerAction which fires out HIDScript on USB host connections is now loaded again and should be shown in the 
TriggerActions tab of the webclient. If P4wnP1 A.L.O.A. is attached to an USB host, the script should run again.

Storing, loading and deploying of templates is one of the two main concepts behind P4wnP1's automation workflow, 
the other one are the already known TriggerActions. It is worth mentioning, that not only TriggerAction sets could be
stored and loaded as templates themselves, but that TriggerActions could be used to deploy already stored templates, if 
that makes sense.

Revisiting our tasks, it seems all defined requirements are met now:
- we typed "Hello world" into the editor of a Windows USB host 
- the editor is opened by P4wnP1, not manually by the user
- the editor is closed automatically, when one of the keyboard LEDs toggled once
- every time P4wnP1 is attached to a USB host, this behavior repeats
- the HIDScript runs only once, unless P4wnP1 is re-attached to the USB host, even if successive keyboard LED changes
occur
- if P4wnP1 is rebooted, the same behavior could be recovered by loading the stored TriggerAction set (which again 
refers to the stored HIDScript). This could either be achieved with a single CLI command or with a simple "load&add" or
"load&replace" from the webclient's trigger action tab.

Once more let us add additional goals:
- it should be assured, that the USB configuration has the keyboard functionality enabled (the current setup doesn't do
this and the TriggerAction couldn't start the HIDScript in case the USB keyboard is disabled)
- the setup shoot applied at boot, without the need to manually load the TriggerAction set. It has to survive a reboot
of P4wnP1.

To achieve the two additional goals, we have to dive into a new topic and ...

#### Introduce Master Templates

Before we look into Master Templates, we do something we haven't done because everything just worked. We define a valid
USB configurations, matching our task:

- device serial number: 123456789
- device product name: Auto Writer
- device manufacturer: The Creator
- Product ID: 0x9876
- Vendor ID: 0x1D6B
- enabled USB functions
  - HID keyboard
  - HID mouse
  
Let's take a look into the usage screen of the proper CLI command first:

``` 
root@kali:~# P4wnP1_cli usb set -h
set USB Gadget settings

Usage:
  P4wnP1_cli usb set [flags]

Flags:
  -e, --cdc-ecm               Use the CDC ECM gadget function
  -n, --disable               If this flag is set, the gadget stays inactive after deployment (not bound to UDC)
  -h, --help                  help for set
  -k, --hid-keyboard          Use the HID KEYBOARD gadget function
  -m, --hid-mouse             Use the HID MOUSE gadget function
  -g, --hid-raw               Use the HID RAW gadget function
  -f, --manufacturer string   Manufacturer string (default "MaMe82")
  -p, --pid string            Product ID (format '0x1347') (default "0x1347")
  -o, --product string        Product name string (default "P4wnP1 by MaMe82")
  -r, --rndis                 Use the RNDIS gadget function
  -s, --serial                Use the SERIAL gadget function
  -x, --sn string             Serial number (alpha numeric) (default "deadbeef1337")
  -u, --ums                   Use the USB Mass Storage gadget function
      --ums-cdrom             If this flag is set, UMS emulates a CD-Rom instead of a flashdrive (ignored, if UMS disabled)
      --ums-file string       Path to the image or block device backing UMS (ignored, if UMS disabled)
  -v, --vid string            Vendor ID (format '0x1d6b') (default "0x1d6b")

Global Flags:
      --host string   The host with the listening P4wnP1 RPC server (default "localhost")
      --json          Output results as JSON if applicable
      --port string   The port on which the P4wnP1 RPC server is listening (default "50051")
``` 

A ton of options, deploying our defined USB setup could be done like this:

```
root@kali:~# P4wnP1_cli usb set \
> --sn 123456789 \
> --product "Auto Writer" \
> --manufacturer "The Creator" \
> --pid "0x9876" \
> --vid "0x1d6b" \
> --hid-keyboard \
> --hid-mouse
Successfully deployed USB gadget settings
Enabled:      true
Product:      Auto Writer
Manufacturer: The Creator
Serialnumber: 123456789
PID:          0x9876
VID:          0x1d6b

Functions:
    RNDIS:        false
    CDC ECM:      false
    Serial:       false
    HID Mouse:    true
    HID Keyboard: true
    HID Generic:  false
    Mass Storage: false
```

The result output of the (long) command shows the intended settings. Let us check the "USB settings" tab of the 
webclient to confirm. All changes should be reflected, if nothing went wrong.

Although it is perfectly possible to deploy a USB setup via CLI, there are several benefits using the webclient in favor
of the CLI in this case:
- changing the settings from the webclient is easier and more convenient
- the webclient has a local state, thus settings could be changed (one by one), without deploying them 
- the webclient could store its current settings (which don't have to be deployed necessarily) to persistent templates
- the CLI client, (currently) isn't able to store USB settings

So it would have been a better choice to use the webclient to change the USB settings to our needs. The good thing
about the CLI approach used here: If the settings have already been deployed, we know that they work (before storing
them to a template).

Again we hit the "store" button, this time in the "USB settings" tab. Once more we call the template `tutorial1` (there
is no conflict with the TriggerAction template stored under the same name, because a different namespace is used for 
USB settings).


We have to stored templates now:
- one for the TriggerAction set, named `tutorial1`
- one for the USB settings, also name `tutorial1`

Assuming the state (of current USB settings, TriggerActions or both) changed somehow, we could reload the settings with 
a single CLI command:

```
P4wnP1_cli template deploy --usb tutorial1 --trigger-actions tutorial1
```

But even if we define a new TriggerAction doing this for the on "service start" trigger, for example with an action 
executing a bash script, we run into a chicken-egg problem.

We can't deploy settings 