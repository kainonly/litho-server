package van.api.entities

import org.hibernate.annotations.ColumnDefault
import org.hibernate.type.BooleanType
import org.hibernate.type.TextType
import javax.persistence.*

@Entity(name = "v_resource")
class Resource {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(columnDefinition = "int(10) unsigned")
    var id: Int? = 0

    @Column(name = "`key`", unique = true, nullable = false, length = 40)
    var key: String? = null

    @Column(nullable = false, columnDefinition = "varchar(40) default 'origin'")
    var parent: String? = null

    @Column(nullable = false)
    var name: TextType? = null

    @Column(nullable = false, columnDefinition = "bit default 0")
    var nav: BooleanType? = null

    @Column(nullable = false, columnDefinition = "bit default 0")
    var router: BooleanType? = null

    @Column(nullable = false, columnDefinition = "bit default 0")
    var policy: BooleanType? = null

    @Column(nullable = true, length = 50)
    var icon: String? = null

    @Column(nullable = false, columnDefinition = "bit default 1")
    var status: BooleanType? = null

    @Column(nullable = false, columnDefinition = "int(10) unsigned default 0")
    var createTime: Int = 0

    @Column(nullable = false, columnDefinition = "int(10) unsigned default 0")
    var updateTime: Int = 0
}

