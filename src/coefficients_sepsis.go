package src

// 1.04 - 1.05*sin( 0.04*(-x_1 - 0.72)**2 - 0.28*sqrt(x_11 + 0.37) + 0.05*log(4.25 - 1.38*x_16) - 0.24*log(3.4*x_14 + 3.95) +
// 	0.12*sin(0.27*x_17 + 1.85) - 0.12*sin(0.31*x_19 + 5.04) + 0.02*sin(0.89*x_20 - 0.18) - 0.49*sin(0.43*x_22 + 2.24) +
// 	0.13*sin(0.41*x_23 + 2.39) + 0.24*sin(0.31*x_24 + 1.61) + 0.35*sin(0.16*x_25 - 4.2) + 0.13*sin(0.18*x_26 - 7.56) +
// 	0.5*sin(0.17*x_27 + 8.5) + 0.13*sin(0.21*x_29 + 2.16) - 0.05*sin(0.4*x_3 + 1.37) - 0.21*sin(0.23*x_30 - 7.04) +
// 	0.18*sin(0.26*x_35 + 2.15) - 0.06*sin(0.28*x_36 - 7.78) + 0.04*sin(0.26*x_37 + 4.52) + 0.05*sin(1.06*x_8 - 9.61) +
// 	0.1*tan(0.28*x_10 - 5.95) + 0.01*tan(0.14*x_4 + 1.0) - 0.07*tanh(0.58*x_2 - 0.48) + 0.03*tanh(0.95*x_21 - 0.53) +
// 	0.02*tanh(1.02*x_28 - 0.68) - 0.11*tanh(0.35*x_31 - 1.43) + 0.04*tanh(1.22*x_32 - 2.35) + 0.03*tanh(1.51*x_33 - 2.12) -
// 	0.32*tanh(0.2*x_5 - 0.85) + 0.02*Abs(9.96*x_18 + 7.21) - 0.01*Abs(6.07*x_34 + 2.42) + 5.76)