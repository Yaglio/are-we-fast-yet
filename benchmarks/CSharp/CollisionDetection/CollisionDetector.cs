using Benchmarks.Som;

namespace Benchmarks.CollisionDetection;

internal sealed class CollisionDetector
{
    private readonly RedBlackTree<CallSign, Vector3D> _state = new();

    public Vector<Collision> HandleNewFrame(Vector<Aircraft> frame)
    {
        Vector<Motion> motions = new();
        RedBlackTree<CallSign, bool> seen = new();

        frame.ForEach(aircraft =>
        {
            Vector3D? oldPosition = _state.Put(aircraft.CallSign, aircraft.Position);
            Vector3D newPosition = aircraft.Position;
            seen.Put(aircraft.CallSign, true);

            oldPosition ??= newPosition;

            motions.Append(new Motion(aircraft.CallSign, oldPosition, newPosition));
        });

        // Remove aircraft that are no longer present.
        Vector<CallSign> toRemove = new();
        _state.ForEach(e =>
        {
            if (!seen.Get(e.Key))
            {
                toRemove.Append(e.Key);
            }
        });

        toRemove.ForEach(e => _state.Remove(e));

        Vector<Vector<Motion>> allReduced = ReduceCollisionSet(motions);
        Vector<Collision> collisions = new();
        allReduced.ForEach(reduced =>
        {
            for (int i = 0; i < reduced.Size(); ++i)
            {
                Motion? motion1 = reduced.At(i);
                for (int j = i + 1; j < reduced.Size(); ++j)
                {
                    Motion? motion2 = reduced.At(j);
                    Vector3D? collision = motion1!.FindIntersection(motion2!);
                    if (collision != null)
                    {
                        collisions.Append(new Collision(motion1.CallSign, motion2!.CallSign, collision));
                    }
                }
            }
        });

        return collisions;
    }

    private static bool IsInVoxel(Vector2D voxel, Motion motion)
    {
        if (voxel.X > Constants.MaxX ||
            voxel.X < Constants.MinX ||
            voxel.Y > Constants.MaxY ||
            voxel.Y < Constants.MinY)
        {
            return false;
        }

        Vector3D init = motion.PosOne;
        Vector3D fin = motion.PosTwo;

        // ReSharper disable once InconsistentNaming
        const double VS = Constants.GoodVoxelSize;
        const double R = Constants.ProximityRadius / 2.0;

        double vX = voxel.X;
        double x0 = init.X;
        double xv = fin.X - init.X;

        double vY = voxel.Y;
        double y0 = init.Y;
        double yv = fin.Y - init.Y;

        double lowX = (vX - R - x0) / xv;
        double highX = (vX + VS + R - x0) / xv;

        if (xv < 0.0)
        {
            double tmp = lowX;
            lowX = highX;
            highX = tmp;
        }

        double lowY = (vY - R - y0) / yv;
        double highY = (vY + VS + R - y0) / yv;

        if (yv < 0.0)
        {
            double tmp = lowY;
            lowY = highY;
            highY = tmp;
        }

        return ((xv == 0.0 && vX <= x0 + R && x0 - R <= vX + VS) /* no motion in x */ ||
                (lowX <= 1.0 && 1.0 <= highX) || (lowX <= 0.0 && 0.0 <= highX) ||
                (0.0 <= lowX && highX <= 1.0)) &&
               ((yv == 0.0 && vY <= y0 + R && y0 - R <= vY + VS) /* no motion in y */ ||
                (lowY <= 1.0 && 1.0 <= highY) || (lowY <= 0.0 && 0.0 <= highY) ||
                (0.0 <= lowY && highY <= 1.0)) &&
               (xv == 0.0 || yv == 0.0 || /* no motion in x or y or both */
                (lowY <= highX && highX <= highY) ||
                (lowY <= lowX && lowX <= highY) ||
                (lowX <= lowY && highY <= highX));
    }

    private static readonly Vector2D _horizontal = new(Constants.GoodVoxelSize, 0.0);
    private static readonly Vector2D _vertical = new(0.0, Constants.GoodVoxelSize);

    private static void PutIntoMap(RedBlackTree<Vector2D, Vector<Motion>> voxelMap, Vector2D voxel, Motion motion)
    {
        Vector<Motion>? array = voxelMap.Get(voxel);
        if (array == null)
        {
            array = new Vector<Motion>();
            voxelMap.Put(voxel, array);
        }

        array.Append(motion);
    }

    private static void Recurse(
        RedBlackTree<Vector2D, Vector<Motion>> voxelMap,
        RedBlackTree<Vector2D, bool> seen,
        Vector2D nextVoxel, Motion motion)
    {
        if (!IsInVoxel(nextVoxel, motion))
        {
            return;
        }

        if (seen.Put(nextVoxel, true))
        {
            return;
        }

        PutIntoMap(voxelMap, nextVoxel, motion);

        Recurse(voxelMap, seen, nextVoxel.Minus(_horizontal), motion);
        Recurse(voxelMap, seen, nextVoxel.Plus(_horizontal), motion);
        Recurse(voxelMap, seen, nextVoxel.Minus(_vertical), motion);
        Recurse(voxelMap, seen, nextVoxel.Plus(_vertical), motion);
        Recurse(voxelMap, seen, nextVoxel.Minus(_horizontal).Minus(_vertical), motion);
        Recurse(voxelMap, seen, nextVoxel.Minus(_horizontal).Plus(_vertical), motion);
        Recurse(voxelMap, seen, nextVoxel.Plus(_horizontal).Minus(_vertical), motion);
        Recurse(voxelMap, seen, nextVoxel.Plus(_horizontal).Plus(_vertical), motion);
    }

    private static Vector<Vector<Motion>> ReduceCollisionSet(Vector<Motion> motions)
    {
        RedBlackTree<Vector2D, Vector<Motion>> voxelMap = new();
        motions.ForEach(motion => DrawMotionOnVoxelMap(voxelMap, motion));

        Vector<Vector<Motion>> result = new();
        voxelMap.ForEach(e =>
        {
            if (e.Value.Size() > 1)
            {
                result.Append(e.Value);
            }
        });
        return result;
    }

    private static Vector2D VoxelHash(Vector3D position)
    {
        int xDiv = (int)(position.X / Constants.GoodVoxelSize);
        int yDiv = (int)(position.Y / Constants.GoodVoxelSize);

        double x = Constants.GoodVoxelSize * xDiv;
        double y = Constants.GoodVoxelSize * yDiv;

        if (position.X < 0)
        {
            x -= Constants.GoodVoxelSize;
        }

        if (position.Y < 0)
        {
            y -= Constants.GoodVoxelSize;
        }

        return new Vector2D(x, y);
    }

    private static void DrawMotionOnVoxelMap(RedBlackTree<Vector2D, Vector<Motion>> voxelMap, Motion motion)
    {
        RedBlackTree<Vector2D, bool> seen = new();
        Recurse(voxelMap, seen, VoxelHash(motion.PosOne), motion);
    }
}

