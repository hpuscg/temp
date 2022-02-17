#pragma once

#include <stdint.h>
#include <stdbool.h>

#define de_crypto_get_uuid      _ZN3std10sys_common17thread_local_dtor22register_dtor_fallback8run_dtor17h8b23aa9fcda6830aE
#define de_crypto_check_auth    _ZN3std10sys_common17thread_local_dtor22register_dtor_fallback9stop_dtor17hb29ed1d869898822E
#define de_crypto_decrypt       _ZN3std10sys_common17thread_local_dtor22register_dtor_fallback9next_dtor17h6b51b4d9673da2eaE
#define de_crypto_lic_import    _ZN3std10sys_common17thread_local_dtor22register_dtor_fallback9iter_dtor17h3c3fc28f7cbd2226E
#define de_crypto_lic_parse     _ZN3std10sys_common17thread_local_dtor22register_dtor_fallback8map_dtor17hf170567088f8423bE
#define de_crypto_lic_elapse    _ZN3std10sys_common17thread_local_dtor22register_dtor_fallback11filter_dtor17h22f55da16a118e0eE

#ifdef __cplusplus
extern "C" {
#endif

/// - 需要的链接参数 -lpthread -ldl
///
/// - 结果指针指向的数据需要由外部释放
///   ``` c
///   uint8_t *dest;
///   int32_t r = de_crypto_xxxxx(&dest);
///   if (r < 0) return r;
///   // do_something(dest, r);
///   free(dest);
///   ```

/// @brief 获取UUID
/// @param dest 结果数据的指针, 指向一个16字节长度的空间
void de_crypto_get_uuid(uint8_t **dest);

/// @brief 检查鉴权
/// @return 鉴权状态
bool de_crypto_check_auth();

/// @brief 数据解密
///     明文数据略小于密文数据长度，根据加密算法和明文数据是否字节对齐不同而略有不同。
///     返回大于 0 的值为明文数据相对于输入地址的偏移量，数据长度为输入长度减去偏移量。
///     该模式下，密文内存由用户申请，解密库 inplace 解密，释放由用户释放前文申请的所有内存。
/// @param data 加密数据块起始地址
/// @param data 加密数据块长度
/// @return
///     >0 解密后数据偏移量
///     -1 输入不合法
///     -2 数据校验失败
///     -3 不支持的加密算法
///     -4 解密失败
///     -254 未鉴权
///     -255 其他错误
int32_t de_crypto_decrypt(uint8_t *data, uint32_t size);

/// 本地证书
/// ``` json
/// {
///     "auth_cnt_left": 3600, // 计数器
///     "auth_days_left": 25,  // 剩余天数
///     "lic_status": true,    // 证书状态
///     "lic_src": {
///         "auth_date_gen": "2019-12-31",      // 计数器
///         "lic_type": "reset_passwd"          // 业务扩展信息
///     }                      // 分发证书原始内容
/// }
/// ```
///
/// @brief 导入证书
///        导入证书后使用check接口获取内容
///        注意分发证书有auth和reset_password两种
/// @param src          分发证书数据
/// @param src_size     分发证书数据长度
/// @param dest         结果数据的指针, 指向分发证书内容
/// @return 是否成功
///     >0 本地证书数据长度
///     -1 数据错误
///     -2 UUID错误
///     -3 签名错误
///     -254 未鉴权
///     -255 失败
int32_t de_crypto_lic_import(uint8_t *src, uint32_t src_size, uint8_t **dest);

/// @brief 检查证书状态
/// @param src          本地证书数据
/// @param src_size     本地证书数据长度
/// @param dest         结果数据的指针, 指向本地证书内容
/// @return 是否成功
///     >0 本地证书内容长度
///     -1 数据错误
///     -254 未鉴权
///     -255 失败
int32_t de_crypto_lic_parse(uint8_t *src, uint32_t src_size, uint8_t **dest);

/// @brief 增加证书ticker
///     每10min调用一次该函数
/// @param src          旧证书数据
/// @param src_size     旧证书数据长度
/// @param dest         结果数据的指针, 指向新证书数据
/// @return 是否成功
///     >0 新证书数据长度
///     -1 数据错误
///     -4 证书已过期
///     -254 未鉴权
///     -255 失败
int32_t de_crypto_lic_elapse(uint8_t *src, uint32_t src_size, uint8_t **dest);

#ifdef __cplusplus
}
#endif