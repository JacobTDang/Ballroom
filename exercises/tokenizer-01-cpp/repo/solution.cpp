#include <sstream>
#include <string>
#include <vector>

struct Token {
    std::string kind;
    std::string text;
    int pos;
};

// Tokenize splits input into tokens; on invalid input, fill *err
// (naming the position) and return {}.
//
// TODO: splitting on spaces calls "3+4" one token and loses every
// position -- and nothing is ever an error.
std::vector<Token> Tokenize(const std::string& input, std::string* err) {
    std::vector<Token> tokens;
    std::istringstream ss(input);
    std::string word;
    while (ss >> word) {
        tokens.push_back({"ident", word, 0});
    }
    return tokens;
}
