namespace Benchmarks.CollisionDetection;

internal sealed class Motion
{
    public readonly CallSign CallSign;
    public readonly Vector3D PosOne;
    public readonly Vector3D PosTwo;

    public Motion(CallSign callSign, Vector3D posOne, Vector3D posTwo)
    {
        CallSign = callSign;
        PosOne = posOne;
        PosTwo = posTwo;
    }

    private Vector3D Delta() => PosTwo.Minus(PosOne);

    public Vector3D? FindIntersection(Motion other)
    {
        Vector3D init1 = PosOne;
        Vector3D init2 = other.PosOne;
        Vector3D vec1 = Delta();
        Vector3D vec2 = other.Delta();
        const double Radius = Constants.ProximityRadius;

        // this test is not geometrical 3-d intersection test, it takes the fact that the aircraft move
        // into account ; so it is more like a 4d test
        // (it assumes that both of the aircraft have a constant speed over the tested interval)

        // we thus have two points, each of them moving on its line segment at constant speed ; we are looking
        // for times when the distance between these two points is smaller than r

        // vec1 is vector of aircraft 1
        // vec2 is vector of aircraft 2

        // a = (V2 - V1)^T * (V2 - V1)
        double a = vec2.Minus(vec1).SquaredMagnitude();

        if (a != 0.0)
        {
            // we are first looking for instances of time when the planes are exactly r from each other
            // at least one plane is moving ; if the planes are moving in parallel, they do not have constant speed

            // if the planes are moving in parallel, then
            //   if the faster starts behind the slower, we can have 2, 1, or 0 solutions
            //   if the faster plane starts in front of the slower, we can have 0 or 1 solutions

            // if the planes are not moving in parallel, then

            // point P1 = I1 + vV1
            // point P2 = I2 + vV2
            //   - looking for v, such that dist(P1,P2) = || P1 - P2 || = r

            // it follows that || P1 - P2 || = sqrt( < P1-P2, P1-P2 > )
            //   0 = -r^2 + < P1 - P2, P1 - P2 >
            //  from properties of dot product
            //   0 = -r^2 + <I1-I2,I1-I2> + v * 2<I1-I2, V1-V2> + v^2 *<V1-V2,V1-V2>
            //   so we calculate a, b, c - and solve the quadratic equation
            //   0 = c + bv + av^2

            // b = 2 * <I1-I2, V1-V2>
            double b = 2.0 * init1.Minus(init2).Dot(vec1.Minus(vec2));

            // c = -r^2 + (I2 - I1)^T * (I2 - I1)
            double c = (-Radius * Radius) + init2.Minus(init1).SquaredMagnitude();

            double discr = (b * b) - (4.0 * a * c);
            if (discr < 0.0)
            {
                return null;
            }

            double v1 = (-b - Math.Sqrt(discr)) / (2.0 * a);
            double v2 = (-b + Math.Sqrt(discr)) / (2.0 * a);

            if (v1 <= v2 && ((v1 <= 1.0 && 1.0 <= v2) ||
                             (v1 <= 0.0 && 0.0 <= v2) ||
                             (0.0 <= v1 && v2 <= 1.0)))
            {
                // Pick a good "time" at which to report the collision.
                double v =
                    // The collision started before this frame. Report it at the start of the frame.
                    // Otherwise, the collision started during this frame. Report it at that moment.
                    v1 <= 0.0 ? 0.0 : v1;

                Vector3D result1 = init1.Plus(vec1.Times(v));
                Vector3D result2 = init2.Plus(vec2.Times(v));

                Vector3D result = result1.Plus(result2).Times(0.5);
                if (result.X is >= Constants.MinX and <= Constants.MaxX &&
                    result.Y is >= Constants.MinY and <= Constants.MaxY &&
                    result.Z is >= Constants.MinZ and <= Constants.MaxZ)
                {
                    return result;
                }
            }

            return null;
        }

        // the planes have the same speeds and are moving in parallel (or they are not moving at all)
        // they  thus have the same distance all the time ; we calculate it from the initial point

        // dist = || i2 - i1 || = sqrt(  ( i2 - i1 )^T * ( i2 - i1 ) )
        double dist = init2.Minus(init1).Magnitude();
        return dist <= Radius ? init1.Plus(init2).Times(0.5) : null;
    }
}
