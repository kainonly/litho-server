package van.api.entities

import org.hibernate.type.TextType
import javax.persistence.*

@Entity(name = "v_admin_basic")
class AdminBasic {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(columnDefinition = "int(10) unsigned")
    var id: Int? = 0

    @Column(nullable = false, length = 255)
    var password: TextType? = null

    @Column(length = 100)
    var email: String? = null

    @Column(length = 20)
    var phone: String? = null

    @Column(length = 30)
    var call: String? = null

    var avatar: String? = null
}