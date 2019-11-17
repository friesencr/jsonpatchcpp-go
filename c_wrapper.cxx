#include "c_wrapper.hxx"
#include "json.hpp"
#include <iostream>

using json = nlohmann::json;

int json_patch(const char* old_json, const char* new_json, char* patch) {
    json j_old = json::parse(old_json);
    json j_new = json::parse(new_json);
    json diffed = json::diff(j_old, j_new);
    // std::string diff_string = diffed.dump();
    return sprintf(patch, "%s", diffed.dump().c_str());
};
