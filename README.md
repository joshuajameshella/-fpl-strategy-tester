# FPL - strategy-tester

Go lockdown project, created to hopefully give me an advantage when picking my next FPL team.

The FPL strategy-tester simulates multiple different strategies available in the Fantasy Premier League.
The Go code used in this repository will help assess each strategy to find the best overall solution. 

All FPL data is taken from [vaastav](https://github.com/vaastav/Fantasy-Premier-League) on Github.


### Results:

The project simulates 10,000 random FPL teams which are then used when running each strategy. The results of which can be found in one of the two following places:

- The overall results and solutions found from simulating each strategy can be found on this [Google Doc](https://docs.google.com/document/d/1NwbvN5KhO3a4yicfFKDgyGPyXOzLO6GolHLPzAUyRaM/edit?usp=sharing). 

- The results (in csv form) can be found in ```internal / simulation_results ```


### Distribution

This strategy determines whether or not the price distribution of players has an effect on the overall points scored during the season.

