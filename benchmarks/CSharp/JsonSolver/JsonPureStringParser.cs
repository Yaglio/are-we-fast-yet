namespace Benchmarks.JsonSolver;

internal sealed class JsonPureStringParser
{
    private readonly string _input;
    private int _index;
    private int _line;
    private int _column;
    private string? _current;
    private string _captureBuffer;
    private int _captureStart;

    public JsonPureStringParser(string str)
    {
        _input = str;
        _index = -1;
        _line = 1;
        _captureStart = -1;
        _column = 0;
        _current = null;
        _captureBuffer = "";
    }

    public JsonValue Parse()
    {
        Read();
        SkipWhiteSpace();
        JsonValue result = ReadValue();
        SkipWhiteSpace();
        if (!IsEndOfText())
        {
            throw Error("Unexpected character");
        }

        return result;
    }

    private JsonValue ReadValue()
    {
        switch (_current)
        {
            case "n":
                return ReadNull();
            case "t":
                return ReadTrue();
            case "f":
                return ReadFalse();
            case "\"":
                return ReadString();
            case "[":
                return ReadArray();
            case "{":
                return ReadObject();
            case "-":
            case "0":
            case "1":
            case "2":
            case "3":
            case "4":
            case "5":
            case "6":
            case "7":
            case "8":
            case "9":
                return ReadNumber();
            default:
                throw Expected("value");
        }
    }

    private JsonArray ReadArray()
    {
        Read();
        JsonArray array = new();
        SkipWhiteSpace();
        if (ReadChar("]"))
        {
            return array;
        }

        do
        {
            SkipWhiteSpace();
            array.Add(ReadValue());
            SkipWhiteSpace();
        } while (ReadChar(","));

        if (!ReadChar("]"))
        {
            throw Expected("',' or ']'");
        }

        return array;
    }

    private JsonObject ReadObject()
    {
        Read();
        JsonObject obj = new();
        SkipWhiteSpace();
        if (ReadChar("}"))
        {
            return obj;
        }

        do
        {
            SkipWhiteSpace();
            string name = ReadName();
            SkipWhiteSpace();
            if (!ReadChar(":"))
            {
                throw Expected("':'");
            }

            SkipWhiteSpace();
            obj.Add(name, ReadValue());
            SkipWhiteSpace();
        } while (ReadChar(","));

        if (!ReadChar("}"))
        {
            throw Expected("',' or '}'");
        }

        return obj;
    }

    private string ReadName()
    {
        if (_current != "\"")
        {
            throw Expected("name");
        }

        return ReadStringInternal();
    }

    private JsonValue ReadNull()
    {
        Read();
        ReadRequiredChar("u");
        ReadRequiredChar("l");
        ReadRequiredChar("l");
        return JsonLiteral.Null;
    }

    private JsonValue ReadTrue()
    {
        Read();
        ReadRequiredChar("r");
        ReadRequiredChar("u");
        ReadRequiredChar("e");
        return JsonLiteral.True;
    }

    private JsonValue ReadFalse()
    {
        Read();
        ReadRequiredChar("a");
        ReadRequiredChar("l");
        ReadRequiredChar("s");
        ReadRequiredChar("e");
        return JsonLiteral.False;
    }

    private void ReadRequiredChar(string ch)
    {
        if (!ReadChar(ch))
        {
            throw Expected("'" + ch + "'");
        }
    }

    private JsonValue ReadString() => new JsonString(ReadStringInternal());

    private string ReadStringInternal()
    {
        Read();
        StartCapture();
        while (_current != "\"")
        {
            if (_current == "\\")
            {
                PauseCapture();
                ReadEscape();
                StartCapture();
            }
            else
            {
                Read();
            }
        }

        string str = EndCapture();
        Read();
        return str;
    }

    private void ReadEscape()
    {
        Read();
        switch (_current)
        {
            case "\"":
            case "/":
            case "\\":
                _captureBuffer += _current;
                break;
            case "b":
                _captureBuffer += "\b";
                break;
            case "f":
                _captureBuffer += "\f";
                break;
            case "n":
                _captureBuffer += "\n";
                break;
            case "r":
                _captureBuffer += "\r";
                break;
            case "t":
                _captureBuffer += "\t";
                break;
            default:
                throw Expected("valid escape sequence");
        }

        Read();
    }

    private JsonValue ReadNumber()
    {
        StartCapture();
        ReadChar("-");
        string? firstDigit = _current;
        if (!ReadDigit())
        {
            throw Expected("digit");
        }

        if (firstDigit != "0")
        {
            while (ReadDigit())
            {
            }
        }

        ReadFraction();
        ReadExponent();
        return new JsonNumber(EndCapture());
    }

    private bool ReadFraction()
    {
        if (!ReadChar("."))
        {
            return false;
        }

        if (!ReadDigit())
        {
            throw Expected("digit");
        }

        while (ReadDigit())
        {
        }

        return true;
    }

    private bool ReadExponent()
    {
        if (!ReadChar("e") && !ReadChar("E"))
        {
            return false;
        }

        if (!ReadChar("+"))
        {
            ReadChar("-");
        }

        if (!ReadDigit())
        {
            throw Expected("digit");
        }

        while (ReadDigit())
        {
        }

        return true;
    }

    private bool ReadChar(string ch)
    {
        if (_current != ch)
        {
            return false;
        }

        Read();
        return true;
    }

    private bool ReadDigit()
    {
        if (!IsDigit())
        {
            return false;
        }

        Read();
        return true;
    }

    private void SkipWhiteSpace()
    {
        while (IsWhiteSpace())
        {
            Read();
        }
    }

    private void Read()
    {
        if (_current == "\n")
        {
            _line++;
            _column = 0;
        }

        _index++;
        if (_index < _input.Length)
        {
            _current = _input.Substring(_index, 1);
        }
        else
        {
            _current = null;
        }
    }

    private void StartCapture()
    {
        _captureStart = _index;
    }

    private void PauseCapture()
    {
        int end = _current == null ? _index : _index - 1;
        _captureBuffer += _input[_captureStart..(end + 1)];
        _captureStart = -1;
    }

    private string EndCapture()
    {
        int end = _current == null ? _index : _index - 1;
        string captured;
        if (_captureBuffer.Length == 0)
        {
            captured = _input[_captureStart..(end + 1)];
        }
        else
        {
            _captureBuffer += _input[_captureStart..(end + 1)];
            captured = _captureBuffer;
            _captureBuffer = "";
        }

        _captureStart = -1;
        return captured;
    }

    private ParseException Expected(string expected)
    {
        if (IsEndOfText())
        {
            return Error("Unexpected end of input");
        }

        return Error($"Expected {expected}");
    }

    private ParseException Error(string message) => new(message, _index, _line, _column - 1);

    private bool IsWhiteSpace() => _current is " " or "\t" or "\n" or "\r";

    private bool IsDigit() => _current is "0" or "1" or "2" or "3" or "4" or "5" or "6" or "7" or "8" or "9";

    private bool IsEndOfText() => _current == null;
}
