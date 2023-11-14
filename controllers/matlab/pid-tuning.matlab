%https://www.mathworks.com/help/control/ref/dynamicsystem.pidtune.html

% Only technical
% 1. Define plant: ramp and sine wave (RMSE), blackbox+least square
% 2. Select Controllers: P, PI, PID, AsTAR, HPA, Others (?), PID+Smooth filter, PID + Delay,
% 3. Tune Controllers: Pid tune: Ziegler-Nichols, AMIGO, Iterative (Matlab)
% 4. Execute experiments, i.e., closed loop with tuned controller
% 5. Evaluate results (fixed goal, variable goal): how to use RMSE? R2? Other metrics?

wc = 1.5;
plant = tf(150.986837,[1 0.079891]); % Ramp plant
%C2 = pidtune(G,'PID2',wc) % original
C2 = pidtune(G,'PID',wc)