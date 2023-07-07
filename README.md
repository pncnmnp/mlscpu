# mlscpu

`mlscpu` is a tool designed for collecting information about the CPU on macOS. This tool takes inspiration from [`lscpu`](https://man7.org/linux/man-pages/man1/lscpu.1.html) tool used in Linux.

## Motivation

Although there are various commands that can retrieve this information (as mentioned in [this Stack Exchange thread](https://apple.stackexchange.com/questions/352769/does-macos-have-a-command-to-retrieve-detailed-cpu-information-like-proc-cpuinf)), there is currently no concrete tool for this purpose.

## Usage

**Note:** Apple Silicon is not yet supported.

To build `mlscpu`, run the following command:

```bash
go build -o mlscpu
```

This will create an executable file named `mlscpu`. Copy this file to the `/usr/local/bin/` directory using the following command:

```bash
cp mlscpu /usr/local/bin/
```

You can now run `mlscpu`
```bash
[14:55][~] mlscpu
Architecture: x86_64
Byte Order: Little Endian
CPU(s): 4
On-line CPU(s): 4
Thread(s) per core: 2
Core(s) per socket: 2
Socket(s): 1
Vendor ID: GenuineIntel
CPU family: 6
CPU Model: 61
Model name: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
Stepping: 4
CPU MHz: 1600
CPU max MHz: 1600
CPU min MHz: 1600
Hyper-Threading Technology: Enabled
L1d cache: 32K
L1i cache: 32K
L2 cache: 256K
L3 cache: 3072K
Flags: FPU VME DE PSE TSC MSR PAE MCE CX8 APIC SEP MTRR PGE MCA CMOV PAT 
       PSE36 CLFSH DS ACPI MMX FXSR SSE SSE2 SS HTT TM PBE SSE3 PCLMULQDQ DTES64 
       MON DSCPL VMX EST TM2 SSSE3 FMA CX16 TPR PDCM SSE4.1 SSE4.2 x2APIC MOVBE 
       POPCNT AES PCID XSAVE OSXSAVE SEGLIM64 TSCTMR AVX1.0 RDRAND F16C
[14:59][mlscpu] 
```
