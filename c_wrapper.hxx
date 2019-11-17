#ifndef C_WRAPPER_HPP
#define C_WRAPPER_HPP


#ifdef __cplusplus

extern "C" {
#endif

int json_patch(const char* oldJosn, const char* newJson, char* patch);

#ifdef __cplusplus
}
#endif

#endif