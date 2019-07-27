# go-traffic-mon
High resolution traffic measurement tool for Linux written in Go.
It allows you to measure interface utilization (bit/s, packets/s, TX & RX) for any interface on a Linux maschine at any interval.
The tool was designed to help localizing very short traffic bursts in networks.

# Usage
```
Usage of ./go-traffic-mon:
  -device string
    	Network device (default "eth0")
  -duration uint
    	Measurement duration in seconds (default 10)
  -rate uint
    	Measurments per second (default 100
```

# Output Format (Example)
```
takt@fuckup:~$ ./traffic -duration=2 -rate=10 -device=eth0
Time (ms)	TX Mbits/s	TX packets/s	RX Mbits/s	RX packets/s
100	2397	71160	32	61870
200	2475	69520	30	59480
300	2368	69920	26	53890
400	2282	70120	29	60690
500	2125	68350	29	57940
600	2551	78870	24	46440
700	2856	82350	30	62060
800	2766	78350	25	48770
900	2537	74710	28	55660
1000	2968	80350	33	63890
1100	3219	88450	36	69680
1200	3298	91980	40	77590
1300	4027	110760	47	92250
1400	4013	114290	40	79700
1500	4553	125770	40	79490
1600	4456	118570	43	85000
1700	4014	107020	34	67490
1800	3619	101120	37	74260
1900	3130	88210	37	75960
```
