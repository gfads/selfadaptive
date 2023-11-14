% File: ExperimentalDesign-BasicPID-RootLocus-1.csv

clear all

maxSize = 102 % maxSize = 15 is the best rmse, r2

% Training - ramp
%yp = [261.574194 621.580645 812.503226 1060.329032 1243.445161 1476.335484 1596.154839 1717.774194 1905.277419 2026.23871 2212.43871 2368.296774 2508.290323 2651.645161 2831.812903 2651.464516 2732.012903 2852.329032 2962.651613 3052.651613 3142.012903 3308.845161 3362.529032 3310.722581 3699.980645 3782.135484 3718.206452 3910.825806 4001.425806 4145.445161 3850.548387 3950.032258 3967.083871 4010.980645 4074.329032 4263.412903 4273.593548 4538.522581 4337.56129 5050.716129 4461.754839 4493.748387 4697.6 4753.664516 5478.554839 4557.96129 4561.567742 4855.148387 4954.245161 4748.103226 4997.567742 5562.832258 4776.006452 5015.63871 4852.503226 4818.412903 4902.025806 5147.432258 5287.477419 5550.283871 5383.877419 5595.845161 5500.632258 5360.877419 5470.941935 5429.767742 5390.6 5480.716129 5424.232258 5670.006452 5594.612903 5556.122581 5613.890323 5500.819355 5632.858065 5837.180645 5767.929032 6035.76129 5768.425806 6019.225806 6292.63871 6253.135484 7040.412903 6242.16129 6230.890323 7166.832258 6611.909677 6558.219355 6395.606452 6367.645161 6373.412903 6494.470968 6641.522581 6421.058065 6575.309677 6597.864516 6485.36129 6522.651613 6238.180645 6498.025806]
%up = [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100]

%Testing - ramp
%yp = [299.922581 587.316129 809.896774 999.174194 1237.929032 1442.470968 1681.683871 1709.96129 1660.154839 1650.393548 1900.974194 2054.729032 2138.812903 2266.270968 2676.193548 2820.709677 2935.2 2993.677419 3303.245161 3345.43871 2998.051613 3232.045161 3833.012903 4014.754839 3980.954839 4196.387097 3491.858065 3892.677419 3727.206452 3726.154839 3853.387097 3938.858065 4078.232258 4110.496774 4179.83871 4237.845161 4254.180645 4212.419355 4115.658065 4247.941935 4199.425806 4319.780645 4626.516129 4576.664516 4511.425806 4706.335484 4697.96129 4865.954839 4879.070968 4890.729032 4778.748387 4901.787097 4816.458065 4899.129032 5156.122581 5451.064516 4971.148387 4862.806452 4541.935484 4921.303226 5342.483871 5111.148387 5355.890323 5252.083871 5262.141935 5261.748387 5133.096774 5376.354839 5852.03871 6131.541935 5963.135484 5475.593548 5550.980645 5490.851613 5455.367742 5432.63871 5242.83871 5576.154839 5595.8 5556.825806 5366.806452 5632.064516 5869.845161 5714.548387 5734.980645 5674.477419 5794.83871 5788.793548 5661.174194 5917.812903 5891.4 5894.451613 5910.019355 5656.974194 5681.354839 6115.032258 5999.058065 6143.787097 6098.722581 5984.916129  ]
%up = [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100]

% Training - Sine wave
yp = [710.122581   266.264516   1770.283871   3771.935484   5044.825806   6027.148387   6478.445161   6373.935484   6072.593548   5351.812903   3706.070968   1958.83871   295.580645   944.651613   2589.554839   3729.458065   4651.625806   5171.309677   5144.6   5035.980645   4507.232258   3672.180645   2264.787097   575.116129   215.935484   1711.748387   3001.522581   4061.832258   4568.812903   5042.103226   4977.767742   4755.367742   4056.116129   2928.451613   1373.341935   229.522581   1024.2   2539.470968   3691.251613   4529.464516   4957.019355   5101.167742   4879.509677   4372.658065   3480.283871   2133.374194   599.032258   227.277419   1847.445161   3463.535484   4187.03871   4785.129032   4877.632258   4913.141935   4482.051613   3913.477419   2605.135484   1248.774194   232.232258   1150.870968   2608.529032   3688.374194   4458.76129   4646.283871   4785.832258   4564.290323   4153.774194   3144.916129   1862.83871   399.722581   414.045161   1870.309677   3207.070968   4073.122581   4601.632258   4842.36129   4671.567742   4422.632258   3546.580645   2541.896774   1096.554839   225.354839   1235.329032   2704.554839   3667.709677   4481.903226   4580.774194   4736.174194   4670.148387   3971.56129   3174.303226   1844.496774   227.483871   409.819355   1998.812903   3256.16129   4267.567742   4784.019355   4804.625806   4783.490323   4310.890323   3825.612903   ]
up = [4 1 11 33 60 84 98 98 84 60 33 11 1 5 22 48 74 93 100 92 72 45 20 3 1 13 35 62 86 99 97 82 57 31 9 1 6 24 50 76 94 100 91 69 43 18 3 1 14 38 65 87 99 97 80 55 28 8 1 7 26 52 78 95 100 89 67 40 16 2 2 16 40 67 89 100 96 79 53 26 7 1 8 28 54 80 96 99 88 65 38 15 1 2 17 42 69 90 100 95 77 51]

%mu = mean(up(1:end-1)); % Original
mu = mean(up(1:maxSize-1));

%my = mean(yp(2:end)); % Original
my = mean(yp(2:maxSize));

u = up - mu;
y = yp - my;

S = zeros(5,1);
%S(1) = sum(y(1:end-1).^2);
%S(2) = sum(u(1:end-1).*y(1:end-1));
%S(3) = sum(u(1:end-1).^2);
%S(4) = sum(y(1:end-1).*y(2:end));
%S(5) = sum(u(1:end-1).*y(2:end));

S(1) = sum(y(1:maxSize-1).^2);
S(2) = sum(u(1:maxSize-1).*y(1:maxSize-1));
S(3) = sum(u(1:maxSize-1).^2);
S(4) = sum(y(1:maxSize-1).*y(2:maxSize));
S(5) = sum(u(1:maxSize-1).*y(2:maxSize));

a = (S(3)*S(4)-S(2)*S(5))/(S(1)*S(3)-(S(2))^2);
b = (S(1)*S(5)-S(2)*S(4))/(S(1)*S(3)-(S(2))^2);

%disp([a b])
%disp([mu my])
%[a b] % a is the root
%[mu my] % operating point

%yhat = a*y(1:end-1) + b*u(1:end-1); % singlestep prediction page 55
yhat = a*y(1:maxSize-1) + b*u(1:maxSize-1);
%yhat2 = a*yhat(1:maxSize-1) + b*u(1:maxSize-1); % multistep prediction page 55

RMSE = rmse(y(2:maxSize),yhat(1:maxSize-1))
R2 = 1-var([y(2:maxSize) yhat(1:maxSize-1)])/var(y)

fprintf('[1 %i];(%f,%f);(%f,%f);%f;%f\n', maxSize,mu,my,a,b, RMSE, R2)

%plot(y(2:end),yhat, '*',y,y,'-');
%plot(y(2:maxSize),yhat(1:maxSize-1), '*',y(1:maxSize-1),y(1:maxSize-1),'-');
%plot(y(2:maxSize),yhat(1:maxSize-1));
plot(up(1:maxSize),yp(1:maxSize),".")
%ylabel("Predicted Rate")
%xlabel("Actual Rate")