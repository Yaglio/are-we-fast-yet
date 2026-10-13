namespace Benchmarks.JsonSolver;

internal class ParseException : Exception
{
    public int Offset { get; }
    public int Line { get; }
    public int Column { get; }

    public ParseException(string message, int offset, int line, int column) : base($"{message} at {line}:{column}")
    {
        Offset = offset;
        Line = line;
        Column = column;
    }
}
