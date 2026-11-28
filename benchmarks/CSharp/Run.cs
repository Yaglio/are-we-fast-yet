using System.Diagnostics;

namespace Benchmarks;

internal sealed class Run(string name)
{
  private int _iterations = 1;
  private int _innerIterations = 1;
  private long _total;
  private readonly string _name = name;
  private readonly Func<Benchmark> _benchmarkSuite = GetSuiteFromName(name);

  private static Func<Benchmark> GetSuiteFromName(string name)
  {
    return name switch
    {
      "Bounce" => static () => new Bounce(),
      "CD" => static () => new CD(),
      "DeltaBlue" => static () => new DeltaBlue(),
      "Havlak" => static () => new Havlak(),
      "Json" => static () => new Json(),
      "List" => static () => new List(),
      "Mandelbrot" => static () => new Mandelbrot(),
      "NBody" => static () => new NBody(),
      "Permute" => static () => new Permute(),
      "Queens" => static () => new Queens(),
      "Richards" => static () => new Richards(),
      "Sieve" => static () => new Sieve(),
      "Storage" => static () => new Storage(),
      "Towers" => static () => new Towers(),
      _ => throw new ArgumentException($"No benchmark found with the name: {name}")
    };
  }

  public void RunBenchmark()
  {
    Console.WriteLine($"Starting {_name} benchmark ...");

    DoRuns(_benchmarkSuite());
    ReportBenchmark();

    Console.WriteLine();
  }

  private void Measure(Benchmark bench)
  {
    long start = Stopwatch.GetTimestamp();
    if (!bench.InnerBenchmarkLoop(_innerIterations))
    {
      throw new InvalidOperationException("Benchmark failed with incorrect result");
    }

    long runTime = (long)(Stopwatch.GetElapsedTime(start).TotalMicroseconds);
    PrintResult(runTime);
    _total += runTime;
  }

  private void DoRuns(Benchmark bench)
  {
    for (int i = 0; i < _iterations; i++)
    {
      Measure(bench);
    }
  }

  private void ReportBenchmark()
  {
    long avgTimeUs = _total / _iterations;
    Console.WriteLine($"{_name}: iterations={_iterations} average: {avgTimeUs}us total: {_total}us\n");
  }

  private void PrintResult(long runTime) => Console.WriteLine($"{_name}: iterations=1 runtime: {runTime}us");

  public void PrintTotal() => Console.WriteLine($"Total Runtime: {_total}us");

  public void SetIterations(int value) => _iterations = value;

  public void SetInnerIterations(int value) => _innerIterations = value;
}
