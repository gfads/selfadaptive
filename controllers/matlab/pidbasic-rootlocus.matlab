% System model: a=0.0865  b=314.4459 y(k+1) = 0.0865 * y(k) + 314.4459 * u(k)

a=0.0865
b=314.4459

%yhat = a*y(1:end-1) + b*u(1:end-1);
%plot(y(2:end),yhat, ' * ',y,y,' -');

%Example: 5*d^2y/dt^2 + 10*dy/dt + 20*y = 30*u
% Define the coefficients
%numerator = 30;
%denominator = [5, 10, 20];

% System: dy/dt + 0.0865*y = 314.4459*u

numerator = 314.4459;
denominator = [1,0.0865];

% Create the transfer function % discrete time
sys = filt(numerator, denominator);

% Design a PID controller
Kp = 1;  % Proportional gain
Ki = 0.5;  % Integral gain
Kd = 0.2;  % Derivative gain

controller = pid(Kp, Ki, Kd);

% Closed-loop system with feedback
closedLoopSys = feedback(controller * sys, 1); %positive feedback loop

% Simulate the response to a step input
time = 0:0.01:10;
u = ones(size(time));  % Step input
[y, t] = lsim(closedLoopSys, u, time);

% Plot the response
figure;
plot(t, y);
xlabel('Time (s)');
ylabel('Output');
title('Step Response of the Closed-Loop System');
