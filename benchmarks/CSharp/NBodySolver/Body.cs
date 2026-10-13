namespace Benchmarks.NBodySolver;

internal sealed class Body
{
    private const double Pi = 3.141592653589793;
    private const double SolarMass = 4 * Pi * Pi;
    private const double DaysPerYer = 365.24;

    public Body(double x, double y, double z, double vx, double vy, double vz, double mass)
    {
        X = x;
        Y = y;
        Z = z;
        Vx = vx * DaysPerYer;
        Vy = vy * DaysPerYer;
        Vz = vz * DaysPerYer;
        Mass = mass * SolarMass;
    }

    public double X { get; set; }
    public double Y { get; set; }
    public double Z { get; set; }
    public double Vx { get; set; }
    public double Vy { get; set; }
    public double Vz { get; set; }
    public double Mass { get; }

    public static Body Jupiter => new(
        4.84143144246472090e+00,
        -1.16032004402742839e+00,
        -1.03622044471123109e-01,
        1.66007664274403694e-03,
        7.69901118419740425e-03,
        -6.90460016972063023e-05,
        9.54791938424326609e-04);

    public static Body Saturn => new(
        8.34336671824457987e+00,
        4.12479856412430479e+00,
        -4.03523417114321381e-01,
        -2.76742510726862411e-03,
        4.99852801234917238e-03,
        2.30417297573763929e-05,
        2.85885980666130812e-04);

    public static Body Uranus => new(
        1.28943695621391310e+01,
        -1.51111514016986312e+01,
        -2.23307578892655734e-01,
        2.96460137564761618e-03,
        2.37847173959480950e-03,
        -2.96589568540237556e-05,
        4.36624404335156298e-05);

    public static Body Neptune => new(
        1.53796971148509165e+01,
        -2.59193146099879641e+01,
        1.79258772950371181e-01,
        2.68067772490389322e-03,
        1.62824170038242295e-03,
        -9.51592254519715870e-05,
        5.15138902046611451e-05);

    public static Body Sun => new(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0);

    internal void OffsetMomentum(double px, double py, double pz)
    {
        Vx = 0.0 - (px / SolarMass);
        Vy = 0.0 - (py / SolarMass);
        Vz = 0.0 - (pz / SolarMass);
    }
}
