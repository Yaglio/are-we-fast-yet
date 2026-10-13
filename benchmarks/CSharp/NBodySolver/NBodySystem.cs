namespace Benchmarks.NBodySolver;

internal sealed class NBodySystem
{
    private readonly Body[] _bodies;

    public NBodySystem()
    {
        _bodies = CreateBodies();
    }

    public Body[] CreateBodies()
    {
        Body[] bodies = [Body.Sun, Body.Jupiter, Body.Saturn, Body.Uranus, Body.Neptune];

        double px = 0.0;
        double py = 0.0;
        double pz = 0.0;

        foreach (Body b in bodies)
        {
            px += b.Vx * b.Mass;
            py += b.Vy * b.Mass;
            pz += b.Vz * b.Mass;
        }

        bodies[0].OffsetMomentum(px, py, pz);

        return bodies;
    }

    public void Advance(double dt)
    {
        for (int i = 0; i < _bodies.Length; i++)
        {
            Body iBody = _bodies[i];
            for (int j = i + 1; j < _bodies.Length; j++)
            {
                Body jBody = _bodies[j];
                double dx = iBody.X - jBody.X;
                double dy = iBody.Y - jBody.Y;
                double dz = iBody.Z - jBody.Z;

                double dSquared = (dx * dx) + (dy * dy) + (dz * dz);
                double distance = Math.Sqrt(dSquared);
                double mag = dt / (dSquared * distance);

                iBody.Vx -= dx * jBody.Mass * mag;
                iBody.Vy -= dy * jBody.Mass * mag;
                iBody.Vz -= dz * jBody.Mass * mag;

                jBody.Vx += dx * iBody.Mass * mag;
                jBody.Vy += dy * iBody.Mass * mag;
                jBody.Vz += dz * iBody.Mass * mag;
            }
        }

        foreach (Body body in _bodies)
        {
            body.X += dt * body.Vx;
            body.Y += dt * body.Vy;
            body.Z += dt * body.Vz;
        }
    }

    public double Energy()
    {
        double e = 0.0;

        for (int i = 0; i < _bodies.Length; i++)
        {
            Body iBody = _bodies[i];
            e += 0.5 * iBody.Mass
                     * ((iBody.Vx * iBody.Vx) +
                        (iBody.Vy * iBody.Vy) +
                        (iBody.Vz * iBody.Vz));

            for (int j = i + 1; j < _bodies.Length; j++)
            {
                Body jBody = _bodies[j];
                double dx = iBody.X - jBody.X;
                double dy = iBody.Y - jBody.Y;
                double dz = iBody.Z - jBody.Z;

                double distance = Math.Sqrt((dx * dx) + (dy * dy) + (dz * dz));
                e -= iBody.Mass * jBody.Mass / distance;
            }
        }

        return e;
    }
}
