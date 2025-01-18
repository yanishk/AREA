import 'package:flutter/material.dart';
import 'package:flutter_switch/flutter_switch.dart';

import "../globalData.dart";
import 'myWorkflowsData.dart';
import '../custom_widget.dart';
import 'myWorkflowsFunctions.dart';

class WorkflowCard extends StatelessWidget {
  const WorkflowCard({
    super.key,
    required this.content,
    this.color = Colors.black,
    this.id = "",
  });

  final Color color;
  final String id;
  final Widget content;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
            decoration: BoxDecoration(
              border: Border.all(
                color: color,
                width: 2.0,
              ),
              borderRadius: BorderRadius.circular(8.0),
            ),
            child: Column(
              children: [
                DecoratedBox(
                  decoration: BoxDecoration(
                    color: color,
                  ),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 25, vertical: 10),
                    child: Padding(
                      padding: EdgeInsets.only(top: 15),
                      child: SizedBox(
                        width: constraints.maxWidth,
                        child: content,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

class GoBackArrow extends StatelessWidget {
  const GoBackArrow({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () {
        Navigator.pop(Globaldata.myContext);
      },
      icon: Icon(
        Icons.arrow_back,
        size: 32,
        color: Colors.white,
      ),
    );
  }
}

class WorkflowsTiles extends StatelessWidget {
  const WorkflowsTiles({
    super.key,
    required this.index,
    required this.onTap,
    this.email = "",
  });
  final int index;
  final void Function() onTap;
  final String email;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Column (
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              SafeWebImage(
                url:
                  "${Globaldata.domainName}${
                    Globaldata.serviceList[MyWorkflowsfunctions.getServiceLocalId(MyWorkflowsData.allWorkflowsData[index].actionID)].icon
                  }",
                width: 40,
                height: 40,
              ),
              Padding(padding: EdgeInsets.only(right: 10),),
              SafeWebImage(
                url:
                  "${Globaldata.domainName}${
                    Globaldata.serviceList[MyWorkflowsfunctions.getServiceLocalId(MyWorkflowsData.allWorkflowsData[index].reactionId)].icon
                  }",
                width: 40,
                height: 40,
              ),
            ],
          ),
          Padding(padding: EdgeInsets.only(top: 20),),
          Text(
            MyWorkflowsData.allWorkflowsData[index].name,
            style: TextStyle(
              fontSize: 25,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          DoubleWidget(
            padding: const EdgeInsets.only(left: 0),
            mainAxisAlignment: MainAxisAlignment.start,
            firstChild: Text(
              "by ",
              style: TextStyle(
                fontSize: 15,
                fontWeight: FontWeight.bold,
                color: const Color.fromARGB(255, 179, 179, 179),
              ),
            ),
            secondChild: Text(
              "$email",
              style: TextStyle(
                fontSize: 15,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
          ),

          Padding(padding: EdgeInsets.only(top: 5),),
          LeftWidget(
            padding: const EdgeInsets.symmetric(horizontal: 0),
            child: FlutterSwitch(
              showOnOff: true,
              inactiveText: "Connect",
              activeText: "Connected",
              toggleColor: HexColor(
                Globaldata.serviceList[MyWorkflowsfunctions.getServiceLocalId(MyWorkflowsData.allWorkflowsData[index].actionID)].color
              ),
              activeTextFontWeight: FontWeight.bold,
              inactiveTextFontWeight: FontWeight.normal,
              inactiveTextColor: Colors.white,
              activeTextColor: const Color.fromARGB(255, 0, 0, 0),
              inactiveColor: const Color.fromARGB(255, 102, 102, 102),
              value: MyWorkflowsData.allWorkflowsData[index].isActivated,
              width: MediaQuery.sizeOf(Globaldata.myContext).width / 3,
              activeColor: const Color.fromARGB(255, 255, 255, 255),
              onToggle: (bool) {},
            ),
          ),
          Padding(padding: EdgeInsets.only(top: 10)),
        ],
      ),
    );
  }
}
