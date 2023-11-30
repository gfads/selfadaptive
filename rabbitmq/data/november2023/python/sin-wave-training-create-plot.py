import matplotlib.pyplot as plt
import numpy as np
import math

o = 50
f = 0.2
a = 50

# Data for plotting
t = np.arange(start=5, stop=1005, step=5)
pc = [1,1,1,3,5,7,11,14,18,23,27,32,38,43,48,54,59,65,70,75,79,84,87,91,94,96,98,99,100,100,99,98,97,94,92,88,84,80,76,71,66,61,55,50,44,39,34,28,24,19,15,11,8,5,3,2,1,1,1,1,2,4,7,10,13,17,21,26,31,36,41,47,52,58,63,68,73,78,82,86,90,93,95,97,99,100,100,100,99,97,95,92,89,86,82,77,72,67,62,57,51,46,40,35,30,25,20,16,12,9,6,4,2,1,1,1,1,2,4,6,9,12,16,20,24,29,34,40,45,51,56,61,67,72,77,81,85,89,92,95,97,99,100,100,100,99,98,96,93,90,87,83,79,74,69,64,58,53,47,42,37,31,26,22,17,13,10,7,4,2,1,1,1,1,1,3,5,8,11,15,19,23,28,33,38,44,49,54,60,65,70,75,80,84,88,91,94,96,98,99]
r = [159.874996,218.682708,231.981361,658.179459,997.130038,1309.764565,1896.503433,2094.849847,2464.401959,2759.510725,2800.747589,3042.140694,3312.949273,3537.123439,3723.721791,4048.351504,4246.763091,4362.715289,4549.500981,4990.784069,4906.127815,4835.855386,4940.200987,4955.807805,5026.45339,5328.030694,4929.116043,4753.756537,4895.514242,4959.27097,5060.42351,4949.751348,4886.366924,4308.81344,4874.082987,4576.797064,4706.150778,4505.940909,4465.749052,4169.160884,4071.293559,4043.872613,3816.478286,3729.299558,3401.168111,3228.717168,2847.9168,2596.235292,2420.239312,2117.698787,1821.806883,1494.314126,1230.452159,867.488329,563.153224,418.456153,213.087071,213.404982,216.613573,213.552656,397.783492,696.581029,1098.798988,1378.518988,1666.158155,1964.460061,2258.971472,2575.12217,2565.183673,2816.524604,3122.470942,3090.631026,3019.344028,3536.684867,3810.523121,4010.854991,3919.028538,4160.713604,4260.095121,4064.696875,4515.330566,4538.905685,4258.447635,4454.441657,4499.262248,4573.818358,4608.631488,4479.465589,4479.217294,4391.623022,4627.543922,4218.113952,3959.832392,3985.311339,4111.138357,4076.666875,3770.894861,3906.112069,3743.380652,3811.380938,3191.96475,3066.384965,2929.22289,2769.702683,2723.749328,2409.997741,2178.408463,1886.733927,1598.105918,1302.09833,966.681214,710.287114,389.910581,213.332196,209.649062,207.43843,212.624766,390.715758,681.130391,959.566182,1285.086663,1460.59912,1732.081206,2183.56462,2410.632597,2575.763315,2684.491574,3015.502075,2954.679059,3376.714286,3162.688714,3622.799004,3887.123946,3776.375527,3913.112799,3871.134391,3949.866975,4118.035292,4470.085099,4175.859218,4345.322737,4415.764763,4292.633207,4365.669708,4034.533323,3902.515475,3916.915315,3930.915893,3990.875588,4171.723581,4163.616072,4074.258108,4039.754752,4052.794868,3679.683155,3548.361053,3280.581679,2922.327449,2980.895339,2791.140748,2821.230003,2604.717256,2342.837559,2063.194629,1804.817773,1610.925377,1301.297667,956.708312,686.270789,381.920108,204.171251,206.532438,218.675605,211.152755,205.246009,567.423954,819.320534,1093.870112,1426.974069,1764.180461,1866.279137,2007.636739,2356.974781,2592.016764,2738.527023,2927.317676,3222.16512,3341.137972,3362.279691,3486.278524,3840.624219,3590.209389,3713.286714,3814.959856,3984.740914,3882.148017,3875.323223,4151.811799,4290.450375,3882.739482]

fig, ax = plt.subplots(2)

ax[0].plot(t,pc,color="black",marker='.', linestyle="none")
ax[0].set(xlabel='time(s)', ylabel='pc',title='')
ax[1].plot(t,r, color="black", marker='.', linestyle="none")
ax[1].set(xlabel='time(s)', ylabel='rate (msg/s)',title='')
ax[0].label_outer()

#ax.grid()

fig.savefig("/content/drive/MyDrive/research/publications/events/2024-ipdps/fig/fig-sin-wave-training.pdf")
plt.show()
