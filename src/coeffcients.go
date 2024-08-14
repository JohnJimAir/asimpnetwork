package src

// 1988.48 * exp(
// 	0.13*(    0.39*sin(7.07*x_2 - 6.21) + 0.09*sin(9.52*x_3 - 8.15)  + 0.21*sin(3.64*x_5 - 0.62) - 0.13*sin(2.24*x_6 + 8.2)   + 0.01*sin(7.85*x_8 + 7.58) + 0.08*tanh(3.77*x_1 - 1.01) + 0.19*tanh(10.0*x_7 - 8.2)   - 0.e-2*Abs(9.96*x_4 - 3.26) + 0.01*Abs(7.94*x_9 - 0.2) - 1)**2 -
// 	0.09*sin( 2.28*(0.33 - x_3)**2      + 1.38*sin(7.4*x_2 + 1.19)   + 1.64*sin(6.44*x_5 - 2.23) + 0.72*sin(6.11*x_6 - 0.73)  - 0.37*sin(5.2*x_7 + 1.18)  - 0.87*sin(4.95*x_8 + 9.62)  + 0.27*tanh(9.6*x_4 - 2.47)   + 0.29*tanh(5.89*x_9 - 2.45) + 5.55) +
// 	2.93*sin( 0.05*(0.24 - x_7)**2      - 0.3*(0.37 - x_8)**3        - 0.61*(0.43 - x_1)**3      - 0.01*sin(5.08*x_2 - 2.22)  - 0.04*sin(6.62*x_3 + 2.99) + 0.05*sin(7.21*x_4 - 5.79)  - 0.e-2*tan(2.2*x_5 - 9.64)   + 0.08*tan(0.28*x_9 + 1.0)   + 0.07*tanh(3.24*x_6 - 2.6) + 4.04) -
// 	0.01*Abs( 21.1*sin(3.89*x_3 - 7.86) + 19.17*sin(3.86*x_4 - 8.02) + 4.29*sin(3.65*x_7 - 1.43) + 12.29*sin(9.79*x_8 + 4.21) + 2.98*tan(1.13*x_1 - 9.75) + 2.38*tan(1.49*x_5 + 2.53)  + 25.59*tanh(3.94*x_2 - 0.58) + 12.94*tanh(10.0*x_6 - 2.6) + 9.55*tanh(7.8*x_9 - 0.84) + 85.59)
// ) - 31.97

// 1.99 - 7.34*tanh(
// 	0.31*(    0.39*sin(7.07*x_2 - 6.21) + 0.09*sin(9.52*x_3 - 8.15)  + 0.21*sin(3.64*x_5 - 0.62) - 0.13*sin(2.24*x_6 + 8.2)   + 0.01*sin(7.85*x_8 + 7.58) + 0.08*tanh(3.77*x_1 - 1.01) + 0.19*tanh(10.0*x_7 - 8.2)   - 0.e-2*Abs(9.96*x_4 - 3.26) + 0.01*Abs(7.94*x_9 - 0.2) - 1)**2 -
// 	0.21*sin( 2.28*(0.33 - x_3)**2      + 1.38*sin(7.4*x_2 + 1.19)   + 1.64*sin(6.44*x_5 - 2.23) + 0.72*sin(6.11*x_6 - 0.73)  - 0.37*sin(5.2*x_7 + 1.18)  - 0.87*sin(4.95*x_8 + 9.62)  + 0.27*tanh(9.6*x_4 - 2.47)   + 0.29*tanh(5.89*x_9 - 2.45) + 5.55) +
// 	7.04*sin( 0.05*(0.24 - x_7)**2      - 0.3*(0.37 - x_8)**3        - 0.61*(0.43 - x_1)**3      - 0.01*sin(5.08*x_2 - 2.22)  - 0.04*sin(6.62*x_3 + 2.99) + 0.05*sin(7.21*x_4 - 5.79)  - 0.e-2*tan(2.2*x_5 - 9.64)   + 0.08*tan(0.28*x_9 + 1.0)   + 0.07*tanh(3.24*x_6 - 2.6) + 4.04) -
// 	0.03*Abs( 21.1*sin(3.89*x_3 - 7.86) + 19.17*sin(3.86*x_4 - 8.02) + 4.29*sin(3.65*x_7 - 1.43) + 12.29*sin(9.79*x_8 + 4.21) + 2.98*tan(1.13*x_1 - 9.75) + 2.38*tan(1.49*x_5 + 2.53)  + 25.59*tanh(3.94*x_2 - 0.58) + 12.94*tanh(10.0*x_6 - 2.6) + 9.55*tanh(7.8*x_9 - 0.84) + 85.59)
//   + 10.72
// )