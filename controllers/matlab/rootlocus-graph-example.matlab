% Chatgpt: how to create a root locus in matlab?

G = tf([1], [1, 3, 2]); % G(s) = 1 / (s^2 + 3s + 2)
H = 1; % Assuming a unity gain controller

rlocus(G*H); % Create the root locus plot for the closed-loop system

grid on;
title('Root Locus Plot');
xlabel('Re');
ylabel('Im');
legend('Pole Locations');

sgrid(0.5, 2);

rlocfind(G*H); % Interactively determine the gain for a specific pole location
