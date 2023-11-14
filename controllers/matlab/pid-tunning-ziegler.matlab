close all, clear, clc
a=0.079891
b=150.986837

% Define your plant model (replace with your actual plant transfer function)
plant_model = tf([b], [1 -a],5.0);

% Define the initial PID controller with zero gains
Kp = 1.0;
Ti = 1.0;
Td = 3.0;
Tf = 5.0;
N = 1
pid_controller = pidstd(Kp,Ti,Td,N, Tf)

% Create the open-loop system
open_loop = pid_controller * plant_model;

% Simulate the open-loop system with a step input
t = 1:5.0:4000;

[y, t] = step(open_loop,t);
%[y, t] = step(plant_model);
plot(t, y);
