# FPL - strategy-tester

The FPL strategy-tester simulates multiple different strategies available in the Fantasy Premier League.
The Go code used in this repository will help assess each strategy to find the best overall solution. 

All FPL data is taken from [vaastav](https://github.com/vaastav/Fantasy-Premier-League) on Github.


### Distribution

This strategy simulates 10,000 possible team layouts, and determines whether or not the price distribution of players has an affect on the overall points scored during the season.

This strategy can be simulated by running the code found in:

```internal / strategy / distribution.go```